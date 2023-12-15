require("dotenv").config();
const celery = require('celery-node');
const amqp = require("amqplib");
const logger = require('./logger')

function generateUuid() {
  return (
    Math.random().toString() +
    Math.random().toString() +
    Math.random().toString()
  );
}

var channel, connection,responseQueue,correlationId; //global variables
var exchange = "daisy1";
async function connectQueue() {
  try{
  connection = await amqp.connect(
    process.env.RABBITMQ || "amqp://localhost:5672",
    (err, connection) => {
      if (err){
        logger.error("Error connecting");
        setTimeout(() => {
          logger.info("Retrying connection in 10 seconds...");
          connectQueue(); // Retry after 10 seconds
        }, 10000);        
      }
    }
  );
  channel = await connection.createChannel((err, connection) => {
    if (err) throw err;
  });

  await channel.assertExchange(exchange, 'direct', {
    durable: true,
  });

  logger.info("Connected to rabbitmq successfully");
  } catch (err){
    logger.error(err);
    logger.info("Retrying connection in 10 seconds...");
    setTimeout(() => {      
      connectQueue(); // Retry after 10 seconds
    }, 10000);     
  }
}




async function sendData(data,service) {
  try {
    const message = JSON.stringify(data);    
    await channel.publish(exchange,service, Buffer.from(message));


    logger.info(` [x] Sent message '${message}' with ID '${service}'`);
  } catch (error) {
    logger.error("Error sending data to the exchange:", error.message);
    throw error;
  }

  // close the channel and connection
  // await channel.close();
  // await connection.close();
}

async function sendMsg(data, service){
  // data to be sent
  try {
    await sendData(data,service); // pass the data to the function we defined
  } catch (error) {
    logger.error(error);
    // throw error; 
  }
};


module.exports = {connectQueue,sendMsg};
