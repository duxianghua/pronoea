<template>
  <el-drawer
    v-if="showDrawer"
    title="Add HTTP Probe"
    :visible.sync="showDrawer"
    direction="rtl"
    :modal="true"
    :show-close="true"
    :wrapper-closable="false"
    size="60%"
    :before-close="onCancel"
    :destroy-on-close="true"
  >
    <div class="app-container">
      <el-form ref="from" :model="form" :rules="formRules" label-width="120px">
        <el-form-item label="Name" prop="metadata.name">
          <el-input v-model="form.metadata.name" />
        </el-form-item>
        <el-form-item label="Labels">
          <Labels v-model="form.metadata.labels" />
        </el-form-item>
        <el-form-item label="Targets" style="margin-top: 15px;">
          <el-container>
            <el-aside width="100px">
              <el-select v-model="form.spec.module.http.method">
                <el-option label="GET" value="GET" />
                <el-option label="POST" value="POST" />
              </el-select>
            </el-aside>

            <div style="padding: 0 0 0 15px; width: 100%;">
              <div v-for="(item,index) in form.spec.targets" :key="index" style="margin: 0 0 0 0;">
                <el-form-item :prop="'spec.targets.'+index" :rules="formRules.targets">
                  <el-row>
                    <el-input
                      v-model="form.spec.targets[index]"
                      placeholder="URL: https://www.google.com/login"
                      class="input-with-select"
                      style="margin-bottom: 10px;"
                      :disabled="target_inpout_edit==index?false:true"
                      @keyup.enter.native="addHost(index)"
                    >
                      <i slot="suffix" class="el-icon-delete-solid" @click="delTarget(index)" />
                    </el-input>
                  </el-row>
                </el-form-item>
              </div>
            </div>

          </el-container>
          <el-row>

            <el-input
              v-if="form.spec.module.http.method === 'POST'"
              v-model="form.spec.module.http.body"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 4}"
              placeholder="Request Body"
            />

          </el-row>
          <el-collapse style="margin-top: 10px;">
            <el-collapse-item title="Advanced Setting" name="1" style="margin-top: 10px;">
              <el-form label-width="240px">
                <el-form-item label="headers">
                  <el-input
                    v-model="form.spec.module.http.headers"
                    type="textarea"
                    :autosize="{ minRows: 1, maxRows: 2}"
                    placeholder="format: {key:value, key2:value}"
                  />
                </el-form-item>
                <el-form-item label="fail_if_body_matches_regexp">
                  <el-input v-model="form.spec.module.http.fail_if_body_matches_regexp" type="textarea" :autosize="{ minRows: 1, maxRows: 2}" placeholder="format: regexp" />
                </el-form-item>
                <el-form-item label="fail_if_body_not_matches_regexp">
                  <el-input v-model="form.spec.module.http.fail_if_body_not_matches_regexp" type="textarea" :autosize="{ minRows: 1, maxRows: 2}" placeholder="format: regexp" />
                </el-form-item>
                <el-form-item label="fail_if_header_matches">
                  <el-input v-model="form.spec.module.http.fail_if_header_matches" type="textarea" :autosize="{ minRows: 1, maxRows: 2}" placeholder="format: regexp" />
                </el-form-item>
                <el-form-item label="fail_if_header_not_matches">
                  <el-input v-model="form.spec.module.http.fail_if_header_not_matches" type="textarea" :autosize="{ minRows: 1, maxRows: 2}" placeholder="format: regexp" />
                </el-form-item>
                <el-form-item label="NoFollowRedirects">
                  <el-switch v-model="form.spec.module.http.no_follow_redirects" />
                </el-form-item>
                <el-form-item label="insecure_skip_verify">
                  <el-switch v-model="form.spec.module.http.tls_config.insecure_skip_verify" />
                </el-form-item>
              </el-form>
            </el-collapse-item>
          </el-collapse>
        </el-form-item>
        <el-form-item label="Timeout">
          <el-input
            v-model="form.spec.module.timeout"
            placeholder="10s,30s,100s"
          />
        </el-form-item>
        <el-form-item label="ContactGroup" prop="spec.contact">
          <el-select v-model="form.spec.contact" filterable placeholder="Please select">
            <el-option
              v-for="item in contactOptions"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button v-if="isEdit" type="primary" @click="onSubmit">Update</el-button>
          <el-button v-if="!isEdit" type="primary" @click="onSubmit">Create</el-button>
          <el-button @click="onCancel">Cancel</el-button>
        </el-form-item>
      </el-form>
    </div>
  </el-drawer>
</template>

<script>
import { CreateProbe, UpdateProbe } from '@/api/probe'
import { ListContactGroup } from '@/api/ContactGroup'
// import labelFrom from '@/components/labels'
import Labels from '@/components/labels/labels'

export default {
  name: 'ProbeForm',
  components: {
    // labelFrom,
    Labels
  },
  props: {
    isActive: {
      type: Boolean,
      default: false
    },
    formData: {
      type: Object,
      required: true
    },
    callBack: {
      type: Function,
      required: false
    }
  },
  data() {
    var checkURL = (rule, value, callback) => {
      if (!value) {
        callback(new Error('asdf'))
      } else {
        let url
        try {
          url = new URL(value)
          callback()
        } catch (_) {
          callback(new Error('URL Format Error: [https||http]://www.google.com/{URI}'))
        }
      }
    }
    return {
      showDrawer: this.isActive,
      isEdit: false,
      contactOptions: [],
      tempLabels: [{ key: '', value: '' }],
      headers: {},
      hosts: '',
      target_inpout_edit: -1,
      formRules: {
        'metadata.name': [
          { required: true, message: 'please enter name', trigger: 'blur' },
          { min: 3, message: 'min size is 3', trigger: 'blur' }
        ],
        'spec.projects': [
          { required: true, message: 'please select project', trigger: 'change' }
        ],
        'targets': [
          { required: true, message: 'please enter url', trigger: ['blur', 'change'] },
          { validator: checkURL, trigger: ['blur', 'change'] }
        ],
        'spec.contact': [
          { required: true, message: 'please select contact group', trigger: ['blur', 'change'] }
        ]
      },
      data: {
        'kind': 'Probe',
        'apiVersion': 'syncbug.io/v1',
        'metadata': {
          'name': '',
          'labels': { 'project': '' }
        },
        'spec': {
          'targets': [''],
          'contact': '',
          'module': {
            'prober': 'http',
            'timeout': '10s',
            'http': {
              'valid_status_codes': [],
              'valid_http_versions': null,
              'preferred_ip_protocol': '',
              'ip_protocol_fallback': false,
              'skip_resolve_phase_with_proxy': false,
              'no_follow_redirects': false,
              'fail_if_ssl': false,
              'fail_if_not_ssl': false,
              'method': 'GET',
              'headers': null,
              'fail_if_body_matches_regexp': null,
              'fail_if_body_not_matches_regexp': null,
              'fail_if_header_matches': null,
              'fail_if_header_not_matches': null,
              'body': '',
              'tls_config': {
                'insecure_skip_verify': 'false'
              },
              'oauth2': {},
              'basic_auth': {},
              'bearer_token': '',
              'compression': '',
              'body_size_limit': '0B'
            }
          }
        }
      },
      form: {},
      options: [
        {
          value: 'devops',
          label: 'devops'
        }
      ]
    }
  },
  computed: {
    labels: {
      get() {
        var _labels = []
        for (var key in this.form.metadata.labels) {
          _labels.push({ key: key, value: this.form.metadata.labels[key] })
        }
        return _labels
      },
      set(val) {
        this.form.metadata.labels[val['key']] = val['value']
      }
    }
  },
  watch: {
    isActive: {
      handler: function(val, oldval) {
        this.showDrawer = val
      }
    },
    showDrawer: {
      handler() {
        this.$emit('update:isActive', this.showDrawer)
      }
    },
    form: {
      handler() {
        this.target_inpout_edit = this.form.spec.targets.length - 1
      },
      deep: true
    },
    formData: {
      handler: function(obj, oldObj) {
        if (Object.keys(obj).length === 0) {
          this.isEdit = false
          this.form = this.data
        } else {
          console.log('data', this.data)
          console.log('obj', obj)
          this.form = { ...this.data, ...JSON.parse(JSON.stringify(obj)) }
          console.log(this.form)
          if ('project' in this.form.metadata.labels) {
            this.form.spec.contact = obj.metadata.labels['project']
          }
          this.isEdit = true
        }
      },
      immediate: true
    }
  },
  created() {
    this.defaultSelect()
    this.initContactGroups()
  },
  methods: {
    initContactGroups() {
      ListContactGroup().then(res => {
        res.items.forEach(item => {
          this.contactOptions.push(item.metadata.name)
        })
      })
    },
    defaultSelect() {
      this.target_inpout_edit = this.form.spec.targets.length - 1
    },
    delTarget(index) {
      if (this.form.spec.targets.length > 1) {
        this.form.spec.targets.splice(index, 1)
      }
    },
    addHost(index) {
      if (this.form.spec.targets[index] !== '') {
        this.form.spec.targets.push('')
        this.target_inpout_edit = index
      }
    },
    onSubmit() {
      this.$refs['from'].validate((valid) => {
        if (valid) {
          this.form.metadata.labels['project'] = this.form.spec.contact
          if (this.isEdit) {
            UpdateProbe(this.form).then(response => {
              this.onCancel()
              this.callBack()
            }).catch(err => {
              console.log(err)
            })
          } else {
            console.log(this.form)
            CreateProbe(this.form).then(response => {
              this.onCancel()
              this.callBack()
            }).catch(err => {
              console.log(err)
            })
          }
        }
      })
    },
    onCancel() {
      this.showDrawer = false
    },
    addField(items) {
      items.push({ key: '', value: '' })
    },
    delField(index, items) {
      items.splice(index, 1)
    }
  }
}
</script>
<!-- @import url("//unpkg.com/element-ui@2.15.10/lib/theme-chalk/index.css"); -->
<style scoped>
.line{
  text-align: center;
}
/* .input-with-select .input-group__prepend {
    background-color: #fff;
  } */
.el-select {
  width: 100%;
}
</style>

