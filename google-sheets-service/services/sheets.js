const { google } = require("googleapis");
const fs = require("fs");
const { OAuth2Client } = require("google-auth-library");
const CLIENT_ID = process.env.CLIENT_ID;
const CLIENT_SECRET = process.env.CLIENT_SECRET;
const REDIRECT_URI = process.env.REDIRECT_URI;
// console.log(CLIENT_ID, CLIENT_SECRET, REDIRECT_URI)
const oAuth2Client = new OAuth2Client(CLIENT_ID, CLIENT_SECRET, REDIRECT_URI);
// const content = fs.readFileSync("token.json", "utf8");
// const credentials = JSON.parse(content);
// let access_token = credentials.access_token;
// let refresh_token = credentials.refresh_token;

// let access_token = process.env.ACCESS_TOKEN;
let refresh_token = process.env.REFRESH_TOKEN;



const mongoose = require("mongoose");
const Sheet = require("../models/sheets");
const Response = require("../models/response");
const { sheets } = require("googleapis/build/src/apis/sheets");

exports.sheets_create_and_add = async (data) => {
    try {
        // const access_token = AUTH_TOKEN || req.session.tokens.access_token;
        try {
            await oAuth2Client.setCredentials({
                access_token: access_token,
            });
        } catch (err) {
            console.log('Refresh token')
            const { tokens } = await oAuth2Client.refreshToken(refresh_token);
            // console.log(tokens);
            await oAuth2Client.setCredentials(tokens);
        }

        // console.log('Got data: ', data);
        let response = data.createdResponse._id;
        let form = data.createdResponse.form._id;
        let title = data.createdResponse.form.title;
        let user = data.createdResponse.user;
        let createdAnswer = data.createdAnswers;
        let answers_array = [];
        createdAnswer.forEach((answer) => {
            let temp = {};
            temp.question_id = answer.question._id;
            temp.question_text = answer.question.text;
            temp.answer_text = answer.text;
            answers_array.push(temp);
        });
        // console.log(answers_array);

        await insert_response(form, response, answers_array);
        await Sheet.findOne({ form: form })
            .exec()
            .then(async (finder) => {
                try {
                    if (finder) {
                        console.log("Sheet exists");
                        let sheet_id = finder.sheet;
                        await add_data_sheet(
                            sheet_id,
                            data.createdResponse.form,
                            answers_array
                        );
                    } else {
                        const spreadsheetTitle = title + "_" + form.slice(-5);
                        let sheet_id = await add_form_sheet(spreadsheetTitle, form);
                        await add_data_sheet(
                            sheet_id,
                            data.createdResponse.form,
                            answers_array
                        );
                    }
                } catch (error) {
                    if (error) console.log('Sheet insertion error: ',error);
                }
            })
            .catch((error) => {
                console.log(error);
            });
    } catch (err) {
        console.log(err);
    }
};

async function add_data_sheet(sheet_id, form, data) {
    try {
        const sheets = google.sheets({ version: "v4", auth: oAuth2Client });
        console.log(data);
        console.log("Adding data");

        let values = [];
        data.forEach((value) => {
            values.push(value.answer_text);
        });

        rows = [values];

        const resource = {
            values: rows,
        };
        const spreadsheetId = sheet_id;
        const valueInputOption = "RAW";
        const range = "Sheet1";

        try {
            const result = await sheets.spreadsheets.values.append({
                spreadsheetId,
                range,
                valueInputOption,
                resource,
            });
            console.log("Data inserted successfully into spreadsheet");
        } catch (err) {
            throw err;
        }
    } catch (err) {
        throw err;
    }
}

async function add_form_sheet(form_title, form_id) {
    const sheets = google.sheets({ version: "v4", auth: oAuth2Client });
    const resource = {
        properties: {
            title: form_title,
        },
    };
    const spreadsheet = await sheets.spreadsheets.create({
        resource,
        fields: "spreadsheetId",
    });
    console.log(`Spreadsheet ID: ${spreadsheet.data.spreadsheetId}`);
    let sheet = new Sheet({
        _id: new mongoose.Types.ObjectId(),
        form: form_id,
        sheet: spreadsheet.data.spreadsheetId,
    });

    sheet
        .save()
        .then(async (res) => {
            return sheet.sheet;
        })
        .catch((err) => {
            throw err;
        });
}

async function insert_response(form, response, answers_array) {
    try {
        // await Response.findOne({ _id: response })
        //     .then((f) => {
        //         if (f) return;
        //     })
        //     .catch((err) => {
        //         throw err;
        //     });
        answers = [];
        answers_array.forEach((answer) => {
            res = new Response({
                response: response,
                form: form,
                question: answer.question_id,
                answer: answer.answer_text,
            });
            answers.push(res);
        });
        Response.insertMany(answers)
            .then((res) => {
                console.log("Answers inserted successfully into the database");
            })
            .catch((err) => {
                throw err;
            });
    } catch (err) {
        throw err;
    }
}
// exports.sheets_write = async (req,res,next)=>{
//   const id = req.params.formId;

// }
