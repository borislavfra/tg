package generator

var (
	responseSchemaTemplate = `import Joi from '@hapi/joi';

export const responseSchema = Joi.object({
    {{ range .}} {{.Base.Name}}: Joi.{{.Type}}().required(), {{ end }}
}).unknown();`

	indexTemplate = `import { JSONRPCRequest, IResponse } from '@mihanizm56/fetch-api';
import { makeRequestConfig, OptionsType } from './make-request-config';

export type ResponseType = IResponse & {
  data: {
    {{ range .Results}} {{.Base.Name}}: {{.Type}};{{ end }}
  };
};

export const get{{.Base.Name}}Request = (values: OptionsType): Promise<IResponse> =>
  new JSONRPCRequest().makeRequest(makeRequestConfig(values));
`
	makeRequestConfigTemplate = `import { responseSchema } from './response-schema';

type TranslateFunctionType = (
  key: string,
  options?: Record<string, any> | null,
) => string;

type RequestParamsType = {
  {{ range .Args}} {{.Base.Name}}: {{.Type}}; {{ end }}
};

export type OptionsType = {
  translateFunction: TranslateFunctionType;
  params: RequestParamsType;
};

const ENDPOINT = '/{{toLowCamel .InterfaceBase.Name}}/{{toLowCamel .Base.Name}}';

export const makeRequestConfig = ({
  translateFunction,
  params,
}: OptionsType) => ({
  endpoint: ENDPOINT,
  translateFunction,
  responseSchema,
  body: { params },
});`
)
