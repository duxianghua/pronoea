import request from '@/utils/request'

export function getList(params) {
  return request({
    url: '/api/v1/probe/',
    method: 'get',
    params
  })
}

export function CreateProbe(data) {
  return request({
    url: '/api/v1/probe/default/' + data.metadata.name,
    method: 'post',
    data
  })
}

export function UpdateProbe(data) {
  return request({
    url: '/api/v1/probe/default/' + data.metadata.name,
    method: 'put',
    data
  })
}



export function deleteProbe(namespace, name) {
  return request({
    url: '/api/v1/probe/' + namespace + "/" + name,
    method: 'delete',
  })
}

export function statusProbe(namespace, name) {
  return request({
    url: '/api/v1/probe/' + namespace + "/" + name +"/"+"status",
    method: 'get',
  })
}