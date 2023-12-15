const http= require('http');
const app = require('./app');
require('dotenv').config();
const logger = require('./services/logger');

const port= process.env.PORT || 3000;

const server = http.createServer(app);

server.listen(port,() =>{
    logger.info(`server listening on port ${port}`);
    logger.info(`Connect using http://localhost:${port}`);
});