import request from '@/utils/request'

export function ListScenarios(params) {
  return request({
    url: '/scenarios/',
    method: 'get',
    params
  })
}

export function GetScenarios(namespace, name) {
  return request({
    url: '/scenarios/' + name,
    method: 'get',
    params: {namespace: namespace}
  })
}

export function CreateScenarios(data) {
  return request({
    url: '/scenarios/' + data.metadata.name,
    method: 'post',
    data
  })
}

export function UpdateScenarios(data) {
  return request({
    url: '/scenarios/' + data.metadata.name,
    method: 'put',
    data
  })
}

export function DeleteScenarios(name, params) {
  return request({
    url: '/scenarios/' + name,
    method: 'delete',
    params
  })
}

export function StatusScenarios(name, params) {
  return request({
    url: '/scenarios/' + name + '/status',
    method: 'get',
    params
  })
}

export function PatchScenarios(name, params) {
  return request({
    url: '/scenarios/' + name,
    method: 'patch',
    params
  })
}
