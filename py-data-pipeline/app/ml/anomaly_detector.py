import numpy as np
from typing import List, Optional
from datetime import datetime, timedelta
from scipy import stats


class AnomalyDetector:
    def __init__(self, z_score_threshold: float = 2.0, iqr_multiplier: float = 1.5):
        self.z_threshold = z_score_threshold
        self.iqr_multiplier = iqr_multiplier

    def detect_location_anomalies(
        self, history: List[dict], current: dict
    ) -> Optional[str]:
        if len(history) < 5:
            return None
        lats = [p["lat"] for p in history[-10:]]
        lons = [p["lon"] for p in history[-10:]]
        speeds = [p.get("speed", 0) for p in history[-10:]]
        cur_speed = current.get("speed", 0)
        speed_z = np.abs(stats.zscore(speeds + [cur_speed])[-1])
        if speed_z > self.z_threshold and cur_speed > 200:
            return "impossible_speed"
        if len(set(zip(lats, lons))) < 2:
            return None
        distances = [self._haversine(lats[i], lons[i], lats[i + 1], lons[i + 1])
                     for i in range(len(lats) - 1)]
        if not distances:
            return None
        q1, q3 = np.percentile(distances, 25), np.percentile(distances, 75)
        iqr = q3 - q1
        upper = q3 + self.iqr_multiplier * iqr
        dist_from_last = self._haversine(
            lats[-1], lons[-1], current["lat"], current["lon"]
        )
        if dist_from_last > upper and dist_from_last > 5000:
            return "location_jump"
        return None

    def detect_route_anomaly(
        self, route_points: List[dict], deviation_threshold_m: float = 100.0
    ) -> Optional[dict]:
        if len(route_points) < 3:
            return None
        for i in range(1, len(route_points) - 1):
            prev, curr, nxt = route_points[i - 1], route_points[i], route_points[i + 1]
            d1 = self._haversine(prev["lat"], prev["lon"], curr["lat"], curr["lon"])
            d2 = self._haversine(curr["lat"], curr["lon"], nxt["lat"], nxt["lon"])
            direct = self._haversine(prev["lat"], prev["lon"], nxt["lat"], nxt["lon"])
            detour_ratio = (d1 + d2) / direct if direct > 0 else 1
            if detour_ratio > 3.0 and d1 > deviation_threshold_m:
                return {
                    "index": i,
                    "lat": curr["lat"],
                    "lon": curr["lon"],
                    "type": "sharp_detour",
                    "detour_ratio": round(detour_ratio, 2),
                }
        return None

    @staticmethod
    def _haversine(lat1: float, lon1: float, lat2: float, lon2: float) -> float:
        R = 6371000
        phi1, phi2 = np.radians(lat1), np.radians(lat2)
        dphi = np.radians(lat2 - lat1)
        dlambda = np.radians(lon2 - lon1)
        a = np.sin(dphi / 2) ** 2 + np.cos(phi1) * np.cos(phi2) * np.sin(dlambda / 2) ** 2
        c = 2 * np.arctan2(np.sqrt(a), np.sqrt(1 - a))
        return R * c
