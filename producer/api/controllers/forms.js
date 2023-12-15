const Form = require("../models/forms");
const Question = require("../models/question");
const Response = require("../models/response");
const Answer = require("../models/answer");
const mongoose = require("mongoose");
const { connectQueue, sendMsg } = require("../../services/connector");
const logger = require("../../services/logger");

exports.create_form = (req, res, next) => {
  description = req.body && req.body.description ? req.body.description : "";
  const form = new Form({
    _id: new mongoose.Types.ObjectId(),
    title: req.body.title,
    description: description,
  });
  // console.log(form);
  var formres;
  form.save().then((result) => {
    formres = result;
    let questions = req.body.questions;
    const formId = result._id;
    let qout = [];
    questions.forEach((q) => {
      let question = new Question({
        _id: new mongoose.Types.ObjectId(),
        form: formId,
        text: q.text,
        type: q.type,
      });
      qout.push(question);
    });

    // console.log(qout);

    Question.insertMany(qout)
      .then((qs_res) => {
        // console.log(qs_res);
        res.status(201).json({
          message: "Form at /forms created successfully",
          createdForm: {
            _id: formres._id,
            title: formres.title,
            description: formres.description,
            createdAt: formres.createdAt,
          },
          createdQs: qs_res,
        });
      })
      .catch((err) => {
        logger.error(err);
        Form.deleteOne({ _id: result._id }).exec();
        res.status(500).json({
          message: "Invalid question",
          error: err,
        });
      });
  });
};

exports.get_form_data = (req, res, next) => {
  const id = req.params.id;
  Form.findById(id)
    .select("title description _id createdAt")
    .exec()
    .then((data) => {
      const response = {
        form: data._id,
        title: data.title,
        description: data.description,
        createdAt: data.createdAt,
      };
      res.status(200).json(response);
    })
    .catch((err) => {
      res.status(500).json({
        error: err,
      });
    });
};

exports.add_response = (req, res, next) => {
  let service = req.headers.service;
  // console.log(service)
  const form_id = req.body.form;
  const answers = req.body.answers;
  const user = req.body.user;
  const metadata = req.body.metadata;

  let res_id;
  Form.findById(form_id)
    .select("_id title")
    .exec()
    .then((data) => {
      let form_title = data.title;
      let user_Response = new Response({
        _id: new mongoose.Types.ObjectId(),
        form: data._id,
        user: user,
      });
      user_Response
        .save()
        .then((res_out) => {
          res_id = res_out._id;
          console.log(res_id);

          let user_answer = [];
          answers.forEach((val) => {
            let answer = new Answer({
              _id: new mongoose.Types.ObjectId(),
              question: val.question,
              response: res_id,
              text: val.text,
            });
            user_answer.push(answer);
          });
          Answer.insertMany(user_answer)
            .then(async (ans_res) => {
              for (let index = 0; index < ans_res.length; index++) {
                await Question.findOne({ "_id": ans_res[index].question }).
                  exec().
                  then((res) => {
                    ans_res[index].question = res;
                  }).catch((err) => {
                    res.status(500).json({
                      message: "Invalid question mappings",
                      error: err,
                    });
                  });
              }

              let payload = {
                createdResponse: {
                  _id: res_id,
                  form: {
                    _id: res_out.form,
                    title: form_title
                  },
                  user: res_out.user,
                  submittedAt: res_out.createdAt,
                },
                createdAnswers: ans_res,
                metadata: metadata
              };

              try {
                sendMsg({ message: payload }, service);
              } catch (err) {
                logger.error(err);
              }


              res.status(201).json({
                message: "Response at /forms created successfully",
                createdResponse: payload.createdResponse,
                createdAnswers: payload.createdAnswers
              });
            })
            .catch((err) => {
              logger.error(err.message);
              res.status(500).json({
                message: "Invalid answers",
                error: err,
              });
            });
        })
        .catch((err) => {
          Response.deleteOne({ _id: res_id }).exec();        
          res.status(500).json({            
            message: "Invalid response",
            error: err,
          });
        });
    });
};

exports.get_question_data = (req, res, next) => {
  const questionId = req.params.id;

  Question.findById(questionId)
    .select("text type _id form")
    .populate("form", "title")
    .exec()
    .then((question) => {
      if (!question) {
        return res.status(404).json({
          message: "Question not found",
        });
      }

      res.status(200).json({
        question: {
          _id: question._id,
          text: question.text,
          type: question.type,
          form: {
            _id: question.form._id,
            title: question.form.title,
          },
        },
      });
    })
    .catch((err) => {
      res.status(500).json({
        error: err,
      });
    });
};
