import request from '@/utils/request'

export function ListContactGroup(params) {
  return request({
    url: '/api/v1/contactgroup/',
    method: 'get',
    params
  })
}

export function GetContactGroup(namespace, name) {
  return request({
    url: '/api/v1/contactgroup/' + name,
    method: 'get'
  })
}

export function CreateContactGroup(data) {
  return request({
    url: '/api/v1/contactgroup/' + data.metadata.name,
    method: 'post',
    data
  })
}

export function UpdateContactGroup(data) {
  return request({
    url: '/api/v1/contactgroup/' + data.metadata.name,
    method: 'put',
    data
  })
}

export function DeleteContactGroup(name) {
  return request({
    url: '/api/v1/contactgroup/' + name,
    method: 'delete'
  })
}

