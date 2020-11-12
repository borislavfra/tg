import {request as adderAddRequest} from './adder/add';
import {request as adderGetUUIDRequest} from './adder/getUUID';
import {request as adderDoNothingRequest} from './adder/doNothing';
import {request as adderBatchedRequest} from './adder/batched';
import {request as todoListAddRequest} from './todoList/add';
import {request as todoListUpdateRequest} from './todoList/update';
import {request as todoListDeleteRequest} from './todoList/delete';
import {request as todoListCreateRequest} from './todoList/create';
import {request as todoListGetRequest} from './todoList/get';


export const APICreator = {
adder: {
AddRequest: adderAddRequest,
GetUUIDRequest: adderGetUUIDRequest,
DoNothingRequest: adderDoNothingRequest,
BatchedRequest: adderBatchedRequest,
},
todoList: {
CreateRequest: todoListCreateRequest,
GetRequest: todoListGetRequest,
AddRequest: todoListAddRequest,
UpdateRequest: todoListUpdateRequest,
DeleteRequest: todoListDeleteRequest,
},
};