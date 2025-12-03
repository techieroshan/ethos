import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const loginTrend = new Trend('login_duration');
const apiTrend = new Trend('api_duration');
const searchTrend = new Trend('search_duration');

// Test configuration
export const options = {
  stages: [
    // Ramp up to 100 users over 2 minutes
    { duration: '2m', target: 100 },
    // Stay at 100 users for 5 minutes
    { duration: '5m', target: 100 },
    // Ramp up to 500 users over 3 minutes
    { duration: '3m', target: 500 },
    // Stay at 500 users for 10 minutes
    { duration: '10m', target: 500 },
    // Ramp down to 0 users over 2 minutes
    { duration: '2m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    http_req_failed: ['rate<0.1'],    // Error rate should be below 10%
    errors: ['rate<0.1'],             // Custom error rate
  },
};

// Base URL
const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

// Test data
const testUsers = [
  { email: 'user1@test.com', password: 'password123' },
  { email: 'user2@test.com', password: 'password123' },
  { email: 'user3@test.com', password: 'password123' },
  { email: 'user4@test.com', password: 'password123' },
  { email: 'user5@test.com', password: 'password123' },
];

// Authentication tokens cache
const authTokens = new Map();

export default function () {
  const userIndex = Math.floor(Math.random() * testUsers.length);
  const testUser = testUsers[userIndex];

  // Login and get token
  const loginStart = new Date().getTime();
  const loginResponse = http.post(`${BASE_URL}/api/v1/auth/login`, {
    email: testUser.email,
    password: testUser.password,
  });
  const loginEnd = new Date().getTime();

  loginTrend.add(loginEnd - loginStart);

  const loginCheck = check(loginResponse, {
    'login status is 200': (r) => r.status === 200,
    'login has token': (r) => r.json().access_token !== undefined,
  });

  if (!loginCheck) {
    errorRate.add(1);
    console.log(`Login failed for ${testUser.email}: ${loginResponse.status} ${loginResponse.body}`);
    return;
  }

  const token = loginResponse.json().access_token;
  authTokens.set(__VU, token);

  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
  };

  // API calls simulation
  const scenarios = [
    // Dashboard access
    () => {
      const start = new Date().getTime();
      const response = http.get(`${BASE_URL}/api/v1/dashboard`, { headers });
      const end = new Date().getTime();
      apiTrend.add(end - start);

      return check(response, {
        'dashboard status is 200': (r) => r.status === 200,
      });
    },

    // Search people
    () => {
      const start = new Date().getTime();
      const response = http.get(`${BASE_URL}/api/v1/people/search?q=test`, { headers });
      const end = new Date().getTime();
      searchTrend.add(end - start);

      return check(response, {
        'search status is 200': (r) => r.status === 200,
      });
    },

    // View feedback
    () => {
      const start = new Date().getTime();
      const response = http.get(`${BASE_URL}/api/v1/feedback/public`, { headers });
      const end = new Date().getTime();
      apiTrend.add(end - start);

      return check(response, {
        'feedback status is 200': (r) => r.status === 200,
      });
    },

    // Create feedback (lower frequency)
    () => {
      if (Math.random() < 0.1) { // Only 10% of users create feedback
        const start = new Date().getTime();
        const response = http.post(`${BASE_URL}/api/v1/feedback`, {
          content: `Load test feedback from ${testUser.email}`,
          rating: Math.floor(Math.random() * 5) + 1,
          recipient_id: 'user-123', // Mock recipient
        }, { headers });
        const end = new Date().getTime();
        apiTrend.add(end - start);

        return check(response, {
          'create feedback status is 201': (r) => r.status === 201,
        });
      }
      return true;
    },
  ];

  // Execute random scenario
  const scenario = scenarios[Math.floor(Math.random() * scenarios.length)];
  const scenarioCheck = scenario();

  if (!scenarioCheck) {
    errorRate.add(1);
  }

  // Random sleep between 1-5 seconds to simulate user behavior
  sleep(Math.random() * 4 + 1);
}

// Setup function - runs before the test starts
export function setup() {
  console.log('Starting performance test setup...');

  // Health check
  const healthResponse = http.get(`${BASE_URL}/health`);
  if (healthResponse.status !== 200) {
    console.error('Health check failed:', healthResponse.status, healthResponse.body);
    throw new Error('Application is not healthy');
  }

  console.log('Health check passed');

  // Pre-warm cache (optional)
  // You could add cache warming logic here

  return {};
}

// Teardown function - runs after the test completes
export function teardown(data) {
  console.log('Performance test completed');

  // Cleanup if needed
  authTokens.clear();
}

// Handle summary - custom summary output
export function handleSummary(data) {
  const summary = {
    'stdout': textSummary(data, { indent: ' ', enableColors: true }),
    'performance-report.json': JSON.stringify(data, null, 2),
  };

  // Save detailed metrics
  const fs = require('fs');
  fs.writeFileSync('performance-metrics.json', JSON.stringify({
    timestamp: new Date().toISOString(),
    metrics: data.metrics,
    thresholds: options.thresholds,
  }, null, 2));

  return summary;
}

function textSummary(data, options) {
  return `
📊 Performance Test Summary
==========================

Test Duration: ${Math.round(data.metrics.iteration_duration.values.avg / 1000)}s
Total Requests: ${data.metrics.http_reqs.values.count}
Failed Requests: ${data.metrics.http_req_failed.values.rate * 100}%

🚀 Response Times:
  Average: ${Math.round(data.metrics.http_req_duration.values.avg)}ms
  95th percentile: ${Math.round(data.metrics.http_req_duration.values['p(95)']}ms
  99th percentile: ${Math.round(data.metrics.http_req_duration.values['p(99)']}ms

🔥 Throughput:
  Requests/second: ${Math.round(data.metrics.http_reqs.values.rate)}

⚠️  Error Rates:
  HTTP errors: ${Math.round(data.metrics.http_req_failed.values.rate * 100)}%
  Custom errors: ${Math.round(data.metrics.errors.values.rate * 100)}%

📈 Custom Metrics:
  Login duration: ${Math.round(data.metrics.login_duration.values.avg)}ms
  API duration: ${Math.round(data.metrics.api_duration.values.avg)}ms
  Search duration: ${Math.round(data.metrics.search_duration.values.avg)}ms

🎯 Thresholds Met:
${Object.entries(options.thresholds).map(([metric, thresholds]) => {
  const metricData = data.metrics[metric];
  if (!metricData) return `  ${metric}: No data`;

  return thresholds.map(threshold => {
    const [op, value] = threshold.split('<');
    const actual = metricData.values[op === 'p(95)' ? 'p(95)' : 'rate'];
    const passed = op === 'rate' ? actual < parseFloat(value) : actual < parseFloat(value);
    return `  ${metric} ${threshold}: ${passed ? '✅' : '❌'} (actual: ${Math.round(actual * (op === 'rate' ? 100 : 1))}${op === 'rate' ? '%' : 'ms'})`;
  }).join('\n');
}).join('\n')}
`;
}
