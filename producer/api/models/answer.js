const mongoose = require("mongoose");

const answerSchema = new mongoose.Schema({
    _id: mongoose.Schema.Types.ObjectId,
    question: {type: mongoose.Schema.Types.ObjectId, ref: "Question",required: true},
    response: {type: mongoose.Schema.Types.ObjectId, ref: "Response",required: true},
    text: {type: String}, 
})

module.exports = mongoose.model("Answer", answerSchema);