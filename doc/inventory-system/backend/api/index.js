// Vercel Serverless Functions 入口文件
const express = require('express');
const cors = require('cors');
const { spawn } = require('child_process');

const app = express();

// 中间件
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// 健康检查
app.get('/api/health', (req, res) => {
  res.json({
    code: 0,
    message: 'success',
    data: {
      status: 'ok',
      timestamp: new Date().toISOString(),
      environment: process.env.NODE_ENV || 'development'
    }
  });
});

// 所有 API 请求转发到 Go 程序
app.all('/api/*', async (req, res) => {
  try {
    // 启动 Go 程序处理请求
    const goProcess = spawn('./go-server', {
      env: {
        ...process.env,
        REQUEST_METHOD: req.method,
        REQUEST_URI: req.originalUrl,
        CONTENT_TYPE: req.get('Content-Type') || '',
        CONTENT_LENGTH: req.get('Content-Length') || '0',
        // 传递请求体
        REQUEST_BODY: JSON.stringify(req.body)
      },
      stdio: ['pipe', 'pipe', 'pipe']
    });

    let stdout = '';
    let stderr = '';

    goProcess.stdout.on('data', (data) => {
      stdout += data.toString();
    });

    goProcess.stderr.on('data', (data) => {
      stderr += data.toString();
    });

    goProcess.on('close', (code) => {
      if (code !== 0) {
        console.error('Go process error:', stderr);
        return res.status(500).json({
          code: -1,
          message: 'Internal server error',
          error: stderr
        });
      }

      try {
        const result = JSON.parse(stdout);
        res.json(result);
      } catch (e) {
        console.error('Failed to parse Go response:', stdout);
        res.status(500).json({
          code: -1,
          message: 'Invalid response format',
          error: stdout
        });
      }
    });

    // 设置超时
    setTimeout(() => {
      goProcess.kill();
      res.status(408).json({
        code: -1,
        message: 'Request timeout'
      });
    }, 8000); // 8秒超时

  } catch (error) {
    console.error('Server error:', error);
    res.status(500).json({
      code: -1,
      message: 'Internal server error',
      error: error.message
    });
  }
});

// 导出 Vercel 函数
module.exports = (req, res) => {
  app(req, res);
};
