const mongoose = require("mongoose");

const questionSchema = new mongoose.Schema({
    _id: mongoose.Schema.Types.ObjectId,
    form: {type: mongoose.Schema.Types.ObjectId, ref: "Form",required: true},
    text: {type: String}, 
    type: {type: String}
})

module.exports = mongoose.model("Question", questionSchema);