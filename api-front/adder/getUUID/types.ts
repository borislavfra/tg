import {TranslateFunction,ExtraValidationCallback,ProgressOptions,CustomSelectorDataType} from '@mihanizm56/fetch-api';
import {AdderGetUUIDParamsType} from '../_types';


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
bodyParams: AdderGetUUIDParamsType;
additionalFetchParams?: FetchParamsType;
};