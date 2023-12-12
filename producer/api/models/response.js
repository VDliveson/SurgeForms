const mongoose = require("mongoose");

const responseSchema = new mongoose.Schema({
    _id: mongoose.Schema.Types.ObjectId,
    form: {type: mongoose.Schema.Types.ObjectId, ref: "Form",required: true},
    user: {type:mongoose.Schema.Types.ObjectId},
    submittedAt: { type: Date, default: Date.now }
})

module.exports = mongoose.model("Response", responseSchema);