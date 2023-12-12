const http= require('http');
const app = require('./app');
require('dotenv').config();

const port= process.env.PORT || 3000;

const server = http.createServer(app);

server.listen(port,() =>{
    console.log(`server listening on port ${port}`);
    console.log(`Connect using http://localhost:${port}`);
});