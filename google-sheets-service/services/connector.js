require("dotenv").config();
const amqp = require("amqplib");
const SheetService = require("./sheets");

var channel, connection;
let id = "sheets";
var exchange = "daisy1";

var messages = [];
async function connectQueue() {
  try {
    connection = await amqp.connect(process.env.RABBITMQ);
    channel = await connection.createChannel();

    await channel.assertExchange(exchange, "direct", {
      durable: true,
    });
    console.log("Connected to rabbitmq exchange " + exchange);

    const queue = "sheets_queue";
    await channel.assertQueue(queue, { durable: false });
    await channel.bindQueue(queue, exchange, 'sheets');

    channel.consume(
      queue,
      async (data) => {
        let content = Buffer.from(data.content).toString();
        console.log(` [x] Received message: ${content} from ID '${id}'`);
        let message = JSON.parse(content);
        messages.push(message);
        // console.log(messages);
        await SheetService.sheets_create_and_add(message.message);
        // let a = await setTimeout(() => {
        // //   let acknowledgment = {
        // //     message: "Microservice1 acknowledged message",
        // //   };
        // //   channel.sendToQueue(
        // //     data.properties.replyTo,
        // //     Buffer.from(JSON.stringify(acknowledgment)),
        // //     {
        // //       correlationId: data.properties.correlationId,
        // //     }
        // //   );

         
        //   console.log("Work completed");
        // }, 3000);
        
        // channel.ack(data);
      },
      {
        noAck: true,
      }
    );
  } catch (error) {
    console.log(error);
  }
}

module.exports = {connectQueue};
