import {METHODS} from '../_methods';
import {AdderAddParamsType,AdderGetUUIDParamsType,AdderDoNothingParamsType,AdderAddResponseType,AdderGetUUIDResponseType,AdderDoNothingResponseType} from '../_types';
import {IResponse,TranslateFunction,ExtraValidationCallback,ProgressOptions,CustomSelectorDataType} from '@mihanizm56/fetch-api';


export type BatchedParamsType = {
method: keyof typeof METHODS;
params: AdderAddParamsType | AdderGetUUIDParamsType | AdderDoNothingParamsType;
};

type FetchParamsType = {
translateFunction?: TranslateFunction;
isErrorTextStraightToOutput?: boolean;
extraValidationCallback?: ExtraValidationCallback;
customTimeout?: number;
abortRequestId?: string;
progressOptions?: ProgressOptions;
customSelectorData?: CustomSelectorDataType;
selectData?: string;
};

export type RequestParamsType = {
bodyParams: Array<BatchedParamsType>;
additionalFetchParams?: FetchParamsType;
};

export type ResponseType = IResponse&{
data: Array<AdderAddResponseType | AdderGetUUIDResponseType | AdderDoNothingResponseType>;
};