const express = require("express");
const session = require('express-session');
const mongoose = require("mongoose");
const morgan = require("morgan");
const {connectQueue} = require('./services/connector')
require('dotenv').config()

const app = express();

// const SESSION_SECRET = process.env.SESSION_SECRET;

// app.use(session({
//     name: "sheets",
//     secret: SESSION_SECRET,
//     resave: false,
//     saveUninitialized: true,
// }));


async function connectDB(){
    try {
    await mongoose.connect(
        "mongodb+srv://vanshajduggal1234:" +
        process.env.MONGO_ATLAS_PW +
        "@service-1.frcd1yg.mongodb.net/?retryWrites=true&w=majority"
      );
      console.log("Connected to MongoDB server");
    } catch(error){
        console.log(error);
    }    
}

connectQueue();
connectDB();

app.use(morgan("dev"));
app.use(express.urlencoded({
  extended: true
}));
app.use(express.json());


// const sheetRoutes = require("./api/routes/sheets");
// const oauthroutes = require("./api/routes/oauth");
// app.use("/api/sheets", sheetRoutes);
// app.use('/api/oauth',oauthroutes);

module.exports = app;