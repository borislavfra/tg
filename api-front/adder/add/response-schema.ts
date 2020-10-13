import {default as Joi} from '@hapi/joi';

export const responseSchema = Joi.object({firstNumber: Joi.number().required(),secondNumber: Joi.number().required()}).unknown();