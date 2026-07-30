export interface User {
  id: string
  username: string
  email: string
  display_name?: string
  avatar_url?: string
  bio?: string
  phone?: string
  role: string
  is_verified: boolean
  created_at: string
  motor_count?: number
  room_count?: number
  total_distance_km?: number
  badge_count?: number
}

export interface Room {
  id: string
  name: string
  description?: string
  cover_url?: string
  category?: string
  is_private: boolean
  member_count: number
  max_members: number
  region?: string
  origin_city?: string
  destination_city?: string
  status: string
  created_by: string
  created_by_name?: string
  created_at: string
}

export interface DashboardStats {
  total_users: number
  active_users_24h: number
  total_rooms: number
  active_rooms: number
  total_routes: number
  total_distance_km: number
  emergency_events_24h: number
  reports_pending: number
  user_growth: { date: string; count: number }[]
  room_activity: { date: string; count: number }[]
  top_rooms: { id: string; name: string; member_count: number }[]
}

export interface CursorPage<T> {
  data: T[]
  next_cursor?: string
  has_more: boolean
  total?: number
}

export interface EmergencyEvent {
  id: string
  room_id: string
  user_id: string
  type: string
  lat: number
  lon: number
  description?: string
  status: string
  created_at: string
}

export interface Report {
  id: string
  reporter_id: string
  reported_id: string
  type: string
  reason: string
  status: string
  created_at: string
}

export interface ModerationAction {
  id: string
  moderator_id: string
  action: string
  target_type: string
  target_id: string
  reason: string
  created_at: string
}
