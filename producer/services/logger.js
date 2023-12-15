const winston = require('winston')
const ecsFormat = require('@elastic/ecs-winston-format')

const customFormat = winston.format.printf(({ level, message }) => {
  return `${level}: ${message}`;
});


const logger = winston.createLogger({
    level: 'debug',
    transports: [
      new winston.transports.Console({
        format: winston.format.combine(
          winston.format.colorize(),
          winston.format.simple()
        ),
        level: 'info'
      }),      
      new winston.transports.File({
        filename: 'logs/log.json',
        level: 'info',
        format: ecsFormat({ convertReqRes: true })
      }),
      new winston.transports.File({
         filename: 'logs/error.json', 
         level: 'error' ,
         format: ecsFormat({ convertReqRes: true })
      }),

    ]
  })

module.exports = logger