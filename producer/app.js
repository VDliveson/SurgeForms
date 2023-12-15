const express = require("express");
const session = require('express-session');
const morgan = require("morgan");
const mongoose = require("mongoose");
require('dotenv').config()
const logger = require('./services/logger');



const {connectQueue,sendMsg} = require('./services/connector')
const app = express();

const formroutes = require("./api/routes/forms");

// const SESSION_SECRET = process.env.SESSION_SECRET;

// app.use(session({
//     name: "sheets",
//     secret: SESSION_SECRET,
//     resave: false,
//     saveUninitialized: true,
// }));

async function connectDB(){
  try{
    mongoose.connect(
      "mongodb+srv://vd:" +
      process.env.MONGO_ATLAS_PW +
      "@sheets-service.z3pyfmw.mongodb.net/?retryWrites=true&w=majority"
    );  
    logger.info("Connected to MongoDB server successfully");
  }catch(err){
    logger.error(err);
  }
}


connectDB();
connectQueue();

class MyStream {
  write(text) {
      logger.info(text.replace(/\n$/, ''));
  }
}
let myStream = new MyStream()
app.use(require("morgan")("tiny", { stream: myStream }));


app.use(express.urlencoded({
  extended: true
}));
app.use(express.json());


app.get('/',(req,res)=>{
  res.status(200).json({
    message: "Producer API",
  })
})

app.get('/msg',(req,res)=>{
  try{
    sendMsg({"message": "Producer API"},"sheets");
    res.status(200).json({"message":"Producer API rabbitmq"});
  }catch{
    res.status(404).json({"message":"rabbitmq unavailable"});
  }
});

app.use("/api/forms", formroutes);


app.use((req, res, next) => {
  const error = new Error("Not found");
  error.status(404);
  next(error);
});

app.use((error, req, res, next) => {
  res.status(error.status || 500);
  res.json({
    message: error.message,
  });
});

module.exports = app;