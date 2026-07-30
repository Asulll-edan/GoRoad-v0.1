from celery import shared_task
from typing import List, Optional
import httpx
from app.config import settings
from app.cache import cache_get_json, cache_set_json


@shared_task(name="tasks.pre_touring_weather")
def pre_touring_weather():
    active_rooms = _get_active_rooms()
    for room in active_rooms:
        lat, lon = room.get("lat"), room.get("lon")
        if lat and lon:
            _check_and_notify(lat, lon, room)


def _get_active_rooms() -> List[dict]:
    return []


def _check_and_notify(lat: float, lon: float, room: dict):
    weather = _fetch_weather(lat, lon)
    if not weather:
        return
    cache_key = f"cache:weather:pre_touring:{room['id']}:notification_sent"
    cached = cache_get_json(cache_key)
    if cached is not None:
        return
    alerts = _check_severe_weather(weather)
    if alerts:
        cache_set_json(cache_key, True, ttl=3600)


def _fetch_weather(lat: float, lon: float) -> Optional[dict]:
    cache_key = f"cache:weather:current:{lat}:{lon}"
    cached = cache_get_json(cache_key)
    if cached:
        return cached
    api_key = settings.OPENWEATHER_API_KEY
    if not api_key:
        return None
    try:
        resp = httpx.get(
            "https://api.openweathermap.org/data/2.5/weather",
            params={"lat": lat, "lon": lon, "appid": api_key, "units": "metric"},
            timeout=10,
        )
        if resp.status_code == 200:
            data = resp.json()
            cache_set_json(cache_key, data, ttl=1800)
            return data
    except httpx.RequestError:
        pass
    return None


def _check_severe_weather(weather: dict) -> List[str]:
    alerts = []
    conditions = weather.get("weather", [])
    for c in conditions:
        main = c.get("main", "")
        if main in ("Thunderstorm", "Tornado", "Hurricane"):
            alerts.append(f"Severe weather: {main}")
    temp = weather.get("main", {}).get("temp", 25)
    if temp > 40:
        alerts.append(f"Extreme heat: {temp}°C")
    elif temp < 5:
        alerts.append(f"Low temperature: {temp}°C")
    wind_speed = weather.get("wind", {}).get("speed", 0)
    if wind_speed > 15:
        alerts.append(f"High wind: {wind_speed} m/s")
    visibility = weather.get("visibility", 10000)
    if visibility < 1000:
        alerts.append("Low visibility")
    return alerts
