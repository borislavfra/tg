
import Joi from '@hapi/joi';
export const SCHEMAS = {
    getUUID: Joi.object({genUUID: Joi.string().required()}),
    doNothing: Joi.object({
        out: Joi.object({
            childThing: Joi.object({
                anything: Joi.object({}).allow(null)}),
            manyThings: Joi.array()
                .items(
                    Joi.object({
                        anything: Joi.object({}).allow(null)
                    }))
                .required().allow(null),
            name: Joi.string().required(),
            thirdName: Joi.string().required(),
            usefulArray: Joi.array()
                .items(
                    Joi.boolean().required()
                ).required().allow(null)
        }),
        tMap: Joi.object().required()}),
    add: Joi.object({sum: Joi.number().required()
    })
};