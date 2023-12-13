require("dotenv").config();
const celery = require('celery-node');
const amqp = require("amqplib");

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
        console.log("Error connecting");
        setTimeout(() => {
          console.log("Retrying connection in 10 seconds...");
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

  console.log("Connected to rabbitmq successfully");
//   responseQueue = await channel.assertQueue("", { exclusive: true });

//   channel.consume(
//     responseQueue.queue,
//     (data) => {
//       if (data.properties.correlationId === correlationId) {
//         console.log(` [x] Received acknowledgment: ${data.content.toString()}`);
//       }
//     },
//     { noAck: true }
//   );    

  } catch (err){
    console.log(err);
    console.log("Retrying connection in 10 seconds...");
    setTimeout(() => {      
      connectQueue(); // Retry after 10 seconds
    }, 10000);     
  }
}




async function sendData(data,service) {
  // send data to queue
  try {
    const message = JSON.stringify(data);    
    await channel.publish(exchange,service, Buffer.from(message));
    // correlationId = generateUuid();
    // await channel.publish(exchange,service, Buffer.from(message), {
    //     correlationId: correlationId,
    //     replyTo: responseQueue.queue,
    //   });    

    console.log(` [x] Sent message '${message}' with ID '${service}'`);
  } catch (error) {
    console.error("Error sending data to the exchange:", error.message);
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
    console.log(error);
    // throw error; 
  }
};


module.exports = {connectQueue,sendMsg};
