import {IResponse} from '@mihanizm56/fetch-api';

export type AdderAddParamsType = {
firstNumber: number;
secondNumber: number;
};

export type AdderGetUUIDParamsType = {
id: number;
};

export type AdderDoNothingParamsType = {
thing: {
childThing: {
anything: Array<{}|null>;
};
manyThings: Array<{}|null>;
name: string;
thirdName: string;
usefulArray: Array<{}|null>;
};
testMap: object;
};

export type AdderAddResponseType = IResponse&{
data: {
sum: number;
};
};

export type AdderGetUUIDResponseType = IResponse&{
data: {
genUUID: string;
};
};

export type AdderDoNothingResponseType = IResponse&{
data: {
tMap: object;
out: {
usefulArray: Array<{}|null>;
childThing: {
anything: Array<{}|null>;
};
manyThings: Array<{}|null>;
name: string;
thirdName: string;
};
};
};
