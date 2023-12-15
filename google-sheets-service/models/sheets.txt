const mongoose = require("mongoose");

const sheetSchema = new mongoose.Schema({
    _id: mongoose.Schema.Types.ObjectId,
    form: mongoose.Schema.Types.ObjectId,
    sheet: mongoose.Schema.Types.String,
    createdAt: {type: Date,default: Date.now()} 
})

module.exports = mongoose.model("Sheet", sheetSchema);