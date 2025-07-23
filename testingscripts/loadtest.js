import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  vus: 10,
  duration: '5s',
};

export default function() {
  // Health check
  let healthRes = http.get('http://localhost:8080/healthz');
  check(healthRes, { "health status is 200": (r) => r.status === 200 });

  // Create short URL
  let payload = JSON.stringify({ url: "https://example.com" });
  let createRes = http.post('http://localhost:8080/create', payload, {
    headers: { 'Content-Type': 'application/json' },
  });
  check(createRes, { "create status is 200": (r) => r.status === 200 });

  let shortcode;
  try {
    shortcode = JSON.parse(createRes.body).shortcode;
  } catch (e) {
    console.error('Failed to parse shortcode:', e);
  }

  if (!shortcode) {
    console.error('Shortcode not found');
    return;
  }

  // Use shortcode for GET request
let getRes = http.get(`http://localhost:8080/${shortcode}`, { redirects: 0 });
  console.log(`GET /${shortcode} status: ${getRes.status}`);
  check(getRes, { "get status is 302": (r) => r.status === 302 });

  // bad request
  payload = JSON.stringify({ url: "example.com" });
  createRes = http.post('http://localhost:8080/create', payload, {
    headers: { 'Content-Type': 'application/json' },
  });
  check(createRes, { "create status is 400 or error": (r) => r.status >= 400 });


  getRes = http.get(`http://localhost:8080/0a`);
  check(getRes, { "get status is 404": (r) => r.status === 404 });


  sleep(1);
}
