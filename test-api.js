// 简单的API测试脚本
const https = require('http');

// 测试变更记录API
const testAPI = () => {
  const options = {
    hostname: 'localhost',
    port: 8088,
    path: '/api/policies/POL_TEST_001/change-records?days=30&page=1&page_size=10',
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      // 这里需要添加实际的JWT token
      'Authorization': 'Bearer YOUR_JWT_TOKEN_HERE'
    }
  };

  const req = https.request(options, (res) => {
    console.log('Status Code:', res.statusCode);
    console.log('Headers:', res.headers);

    let data = '';
    res.on('data', (chunk) => {
      data += chunk;
    });

    res.on('end', () => {
      try {
        const response = JSON.parse(data);
        console.log('Response:');
        console.log(JSON.stringify(response, null, 2));
      } catch (error) {
        console.log('Raw Response:', data);
      }
    });
  });

  req.on('error', (error) => {
    console.error('Request error:', error);
  });

  req.end();
};

// 测试健康检查
const testHealth = () => {
  const options = {
    hostname: 'localhost',
    port: 8088,
    path: '/health',
    method: 'GET'
  };

  console.log('Testing health endpoint...');
  const req = https.request(options, (res) => {
    console.log('Health Status Code:', res.statusCode);
    
    let data = '';
    res.on('data', (chunk) => {
      data += chunk;
    });

    res.on('end', () => {
      console.log('Health Response:', data);
      
      if (res.statusCode === 200) {
        console.log('✅ Backend is running');
        console.log('\nTesting change records API...');
        testAPI();
      } else {
        console.log('❌ Backend health check failed');
      }
    });
  });

  req.on('error', (error) => {
    console.error('❌ Cannot connect to backend:', error.message);
    console.log('Please ensure the backend is running on port 8088');
  });

  req.end();
};

testHealth(); 