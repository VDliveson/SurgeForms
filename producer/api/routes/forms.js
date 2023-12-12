const express = require("express");
const router = express.Router();
const mongoose = require("mongoose");

const morgan = require("morgan");
require('dotenv').config()
const formController = require("../controllers/forms");

router.post("/create",formController.create_form);
router.get("/get/:id", formController.get_form_data);
router.post("/response", formController.add_response);
router.get("/question/:id", formController.get_question_data);

router.get("/", (req, res) => {
    res.status(200).json({
        "message": "Home for forms api"
    })
});

module.exports = router;