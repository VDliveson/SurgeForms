const mongoose = require("mongoose");

const responseSchema = new mongoose.Schema({
    _id: mongoose.Schema.Types.ObjectId,
    response: {type:mongoose.Schema.Types.ObjectId,required:true},
    form: {type:mongoose.Schema.Types.ObjectId,required:true},
    question: {type:mongoose.Schema.Types.ObjectId,required:true},
    text: mongoose.Schema.Types.String,    
    createdAt: {type: Date,default: Date.now()} 
})

module.exports = mongoose.model("Response", responseSchema);