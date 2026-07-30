import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const failures = new Rate('request_failures');
const latency = new Trend('request_latency_ms');

export const options = {
  stages: [
    { duration: '30s', target: 10 },
    { duration: '1m', target: 50 },
    { duration: '30s', target: 100 },
    { duration: '1m', target: 100 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'],
    request_failures: ['rate<0.05'],
  },
};

function randomId() {
  return crypto.randomUUID();
}

function authFlow() {
  const email = `test_${randomId().slice(0,8)}@example.com`;
  const payload = JSON.stringify({ username: email.split('@')[0], email, password: 'Test12345!' });
  const reg = http.post(`${BASE_URL}/v1/auth/register`, payload, { headers: { 'Content-Type': 'application/json' } });
  check(reg, { 'register status 201': (r) => r.status === 201 || r.status === 409 });
  if (reg.status === 201 || reg.status === 200) {
    const login = http.post(`${BASE_URL}/v1/auth/login`, JSON.stringify({ email, password: 'Test12345!' }), {
      headers: { 'Content-Type': 'application/json' },
    });
    check(login, { 'login status 200': (r) => r.status === 200 });
    if (login.status === 200) {
      return login.json('tokens.access_token');
    }
  }
  return null;
}

function getRooms(token) {
  const res = http.get(`${BASE_URL}/v1/rooms?limit=20`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  check(res, { 'get rooms status 200': (r) => r.status === 200 });
  latency.add(res.timings.duration);
  failures.add(res.status !== 200);
  return res;
}

function createRoom(token) {
  const payload = JSON.stringify({ name: `K6 Room ${randomId().slice(0,6)}`, max_members: 10 });
  const res = http.post(`${BASE_URL}/v1/rooms`, payload, {
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
  });
  check(res, { 'create room status 201': (r) => r.status === 201 });
  return res;
}

function getRoutes(token) {
  const res = http.get(`${BASE_URL}/v1/routes?limit=10`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  check(res, { 'get routes status 200': (r) => r.status === 200 });
  return res;
}

export default function () {
  const token = authFlow();
  if (!token) {
    failures.add(1);
    sleep(1);
    return;
  }
  getRooms(token);
  sleep(0.5);
  createRoom(token);
  sleep(0.5);
  getRoutes(token);
  sleep(Math.random() * 2 + 1);
}
