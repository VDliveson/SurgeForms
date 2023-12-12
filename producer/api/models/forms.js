const mongoose = require("mongoose");

const formSchema = new mongoose.Schema({
    _id: mongoose.Schema.Types.ObjectId,
    title: {type: String,required: true},
    description: {type: String},
    createdAt: {type: Date,default: Date.now()} 
})

module.exports = mongoose.model("Form", formSchema);