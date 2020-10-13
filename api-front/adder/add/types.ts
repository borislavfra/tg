import {IResponse,TranslateFunction,ExtraValidationCallback,ProgressOptions,CustomSelectorDataType} from '@mihanizm56/fetch-api';

type ParamsType = {
firstNumber: number;
secondNumber: number;
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
bodyParams: ParamsType;
additionalFetchParams?: FetchParamsType;
};

export type ResponseType = IResponse&{
data: {
sum: number;
};
};