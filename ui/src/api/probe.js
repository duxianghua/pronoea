import request from '@/utils/request'

export function getList(params) {
  return request({
    url: '/probe/',
    method: 'get',
    params
  })
}

export function CreateProbe(data) {
  return request({
    url: '/probe/' + data.metadata.name,
    method: 'post',
    data
  })
}

export function UpdateProbe(data) {
  return request({
    url: '/probe/' + data.metadata.name,
    method: 'put',
    data
  })
}

export function deleteProbe(name, params) {
  return request({
    url: '/probe/' + name,
    method: 'delete',
    params
  })
}

export function statusProbe(name, params) {
  return request({
    url: '/probe/' + name + '/status',
    method: 'get',
    params
  })
}