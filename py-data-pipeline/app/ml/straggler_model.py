import numpy as np
from typing import List, Optional
from datetime import datetime, timedelta


class StragglerDetector:
    def __init__(
        self,
        speed_threshold_kmh: float = 10.0,
        distance_threshold_m: float = 200.0,
        time_threshold_min: float = 5.0,
    ):
        self.speed_threshold = speed_threshold_kmh
        self.distance_threshold = distance_threshold_m
        self.time_threshold = timedelta(minutes=time_threshold_min)

    def detect(
        self,
        formation_positions: List[dict],
        current_positions: List[dict],
    ) -> List[dict]:
        if not formation_positions or not current_positions:
            return []
        stragglers = []
        leader_pos = self._get_lead_position(formation_positions)
        if not leader_pos:
            return []
        for member in current_positions:
            lat = member.get("lat") or member.get("latitude", 0)
            lon = member.get("lon") or member.get("longitude", 0)
            speed = member.get("speed", 0)
            distance = self._haversine(
                leader_pos["lat"], leader_pos["lon"], lat, lon
            )
            is_slow = speed < self.speed_threshold
            is_far = distance > self.distance_threshold
            if is_slow or is_far:
                stragglers.append(
                    {
                        "user_id": member.get("user_id"),
                        "lat": lat,
                        "lon": lon,
                        "speed": speed,
                        "distance_from_lead_m": round(distance, 1),
                        "reason": "slow" if is_slow else "far",
                        "detected_at": datetime.utcnow().isoformat(),
                    }
                )
        return stragglers

    def _get_lead_position(self, positions: List[dict]) -> Optional[dict]:
        for p in positions:
            if p.get("role") == "lead":
                return {"lat": p["lat"], "lon": p["lon"]}
        return positions[0] if positions else None

    @staticmethod
    def _haversine(lat1: float, lon1: float, lat2: float, lon2: float) -> float:
        R = 6371000
        phi1, phi2 = np.radians(lat1), np.radians(lat2)
        dphi = np.radians(lat2 - lat1)
        dlambda = np.radians(lon2 - lon1)
        a = np.sin(dphi / 2) ** 2 + np.cos(phi1) * np.cos(phi2) * np.sin(dlambda / 2) ** 2
        c = 2 * np.arctan2(np.sqrt(a), np.sqrt(1 - a))
        return R * c
