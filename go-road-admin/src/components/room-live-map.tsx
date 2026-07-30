'use client'

import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet'
import 'leaflet/dist/leaflet.css'
import L from 'leaflet'

// Fix Leaflet default marker icon
delete (L.Icon.Default.prototype as any)._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
})

export default function RoomLiveMap() {
  return (
    <div className="card">
      <h3 className="font-semibold mb-4">Live Map</h3>
      <div className="h-[400px] rounded-lg overflow-hidden">
        <MapContainer center={[-6.2, 106.8]} zoom={12} className="h-full w-full">
          <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
          <Marker position={[-6.2, 106.8]}>
            <Popup>Rider positions appear here in real-time</Popup>
          </Marker>
        </MapContainer>
      </div>
    </div>
  )
}
