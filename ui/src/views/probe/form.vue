<template>
  <el-drawer
    title="Add HTTP Probe"
    :visible.sync="showDrawer"
    v-if="showDrawer"
    direction="rtl"
    :modal="true"
    :show-close="true"
    :wrapperClosable="false"
    size="60%"
    :before-close="onCancel"
    :destroy-on-close="true"
  >
    <div class="app-container">
      <el-form ref="from" :model="form" :rules="formRules" label-width="120px">
        <el-form-item label="Name" prop="metadata.name">
          <el-input v-model="form.metadata.name" />
        </el-form-item>
        <!-- <el-form-item label="Labels">
          <labelFrom :data="tempLabels">
          </labelFrom>
        </el-form-item> -->
        <el-form-item label="Targets" style="margin-top: 15px;">
            <el-container>
              <el-aside width="130px">
                <el-select v-model="form.spec.module.http.method">
                  <el-option label="GET" value="GET"></el-option>
                  <el-option label="POST" value="POST"></el-option>
                </el-select>
            </el-aside>
            <!-- <el-main style="padding: 0 0 0 15px;">
              <div style="margin: 0 0 0 0;" v-for="(item,index) in form.spec.targets" :key="index">
                <el-form-item :prop="'spec.targets.'+index" :rules="formRules.targets">
                  <el-row >
                    <el-input
                      placeholder="URL: https://www.google.com/login"
                      v-model="form.spec.targets[index]"
                      class="input-with-select" 
                      style="margin-bottom: 10px;"
                      @keyup.enter.native="addHost(index)"
                      :disabled="target_inpout_edit==index?false:true"
                    >
                      <i slot="suffix" class="el-icon-delete-solid" @click="delTarget(index)"></i>
                    </el-input>
                  </el-row>
                  </el-form-item>
              </div>
            </el-main> -->
            <div style="padding: 0 0 0 15px; width: 100%;">
              <div style="margin: 0 0 0 0;" v-for="(item,index) in form.spec.targets" :key="index">
                <el-form-item :prop="'spec.targets.'+index" :rules="formRules.targets">
                  <el-row >
                    <el-input
                      placeholder="URL: https://www.google.com/login"
                      v-model="form.spec.targets[index]"
                      class="input-with-select" 
                      style="margin-bottom: 10px;"
                      @keyup.enter.native="addHost(index)"
                      :disabled="target_inpout_edit==index?false:true"
                    >
                      <i slot="suffix" class="el-icon-delete-solid" @click="delTarget(index)"></i>
                    </el-input>
                  </el-row>
                  </el-form-item>
              </div>
            </div>
            </el-container>
            <el-collapse style="margin-top: 50px;">
              <el-collapse-item title="Advanced Setting" name="1" style="margin-top: 10px;">
                <el-form label-width="240px">
                  <el-form-item label="headers">
                    <el-input
                    type="textarea"
                    :autosize="{ minRows: 1, maxRows: 2}"
                    placeholder="format: {key:value, key2:value}"
                    v-model="form.spec.module.http.headers">
                </el-input>
                  </el-form-item>
                  <el-form-item label="fail_if_body_matches_regexp">
                    <el-input type="textarea" :autosize="{ minRows: 1, maxRows: 2}" placeholder="format: regexp" v-model="form.spec.module.http.fail_if_body_matches_regexp"></el-input>
                  </el-form-item>
                  <el-form-item label="fail_if_body_not_matches_regexp">
                    <el-input type="textarea" :autosize="{ minRows: 1, maxRows: 2}" placeholder="format: regexp" v-model="form.spec.module.http.fail_if_body_not_matches_regexp"></el-input>
                  </el-form-item>
                  <el-form-item label="fail_if_header_matches">
                    <el-input type="textarea" :autosize="{ minRows: 1, maxRows: 2}" placeholder="format: regexp" v-model="form.spec.module.http.fail_if_header_matches"></el-input>
                  </el-form-item>
                  <el-form-item label="fail_if_header_not_matches">
                    <el-input type="textarea" :autosize="{ minRows: 1, maxRows: 2}" placeholder="format: regexp" v-model="form.spec.module.http.fail_if_header_not_matches"></el-input>
                  </el-form-item>
                  <el-form-item label="NoFollowRedirects">
                    <el-switch v-model="form.spec.module.http.no_follow_redirects" />
                  </el-form-item>
                </el-form>
              </el-collapse-item>
            </el-collapse>
        </el-form-item>
        <el-form-item label="ContactGroup" prop="spec.contact">
          <el-select v-model="form.spec.contact" filterable placeholder="Please select">
          <el-option
            v-for="item in contactOptions"
            :key="item"
            :label="item"
            :value="item">
          </el-option>
        </el-select>
      </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="onSubmit" v-if="isEdit">Update</el-button>
          <el-button type="primary" @click="onSubmit" v-if="!isEdit">Create</el-button>
          <el-button @click="onCancel">Cancel</el-button>
        </el-form-item>
      </el-form>
    </div>
  </el-drawer>
</template>



<script>
import { CreateProbe, UpdateProbe } from '@/api/probe'
import { ListContactGroup } from '@/api/ContactGroup'
import labelFrom from '@/components/labels'

export default {
  name: "probeForm",
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
    var checkURL = (rule, value, callback)=>{
      if (!value){
        callback(new Error("asdf"))
      }else{
        let url;
        try {
          url = new URL(value);
          callback()
        } catch (_) {
          callback(new Error("URL Format Error: [https||http]://www.google.com/{URI}"))
        }
      }
    };
    return {
      showDrawer: this.isActive,
      isEdit: false,
      contactOptions: [],
      tempLabels:[{key:"project",value:""}],
      headers: {},
      hosts: "",
      target_inpout_edit: -1,
      formRules: {
        'metadata.name': [
          {required: true, message: "please enter name", trigger: "blur"},
          {min: 3, message: "min size is 3", trigger:'blur'}
        ],
        'spec.projects': [
          {required: true, message: "please select project", trigger: "change"},          
        ],
        'targets': [
          {required: true, message: "please enter url", trigger: ["blur", "change"]},
          {validator: checkURL, trigger: ["blur", "change"]}
        ],
        'spec.contact': [
          {required: true, message: "please select contact group", trigger: ["blur", "change"]},
        ]
      },
      data: {
        "kind": "Probe",
        "apiVersion": "syncbug.io/v1",
        "metadata": {
            "name": "",
            "labels": {"project": ""}
        },
      "spec": {
        "targets": [""],
        "contact": "",
        "module": {
            "prober": "http",
            "http": {
                "valid_status_codes": [],
                "valid_http_versions": null,
                "preferred_ip_protocol": "",
                "ip_protocol_fallback": false,
                "skip_resolve_phase_with_proxy": false,
                "no_follow_redirects": false,
                "fail_if_ssl": false,
                "fail_if_not_ssl": false,
                "method": "GET",
                "headers": null,
                "fail_if_body_matches_regexp": null,
                "fail_if_body_not_matches_regexp": null,
                "fail_if_header_matches": null,
                "fail_if_header_not_matches": null,
                "body": "",
                "http_client_config": {
                    "tls_config": {}
                },
                "compression": "",
                "body_size_limit": "0B"
            }
        }
      },
      },
      form: {},
      options: [
      {
          value: 'devops',
          label: 'devops'
        }, {
          value: 'sso',
          label: 'sso'
        }, {
          value: 'cmt',
          label: 'cmt'
        }
      ]
    }
  },
  watch: {
    isActive: {
      handler: function(val, oldval){
        this.showDrawer = val
      }
    },
    showDrawer:{
      handler(){
        this.$emit("update:isActive", this.showDrawer)
      }
    },
    form: {
        handler(){
          this.target_inpout_edit = this.form.spec.targets.length -1
     },
     deep: true
    },
    formData: {
      handler: function(obj, oldObj){
        if (Object.keys(obj).length === 0){
          this.isEdit = false
          this.form = this.data
        }else{
          this.form = JSON.parse(JSON.stringify(obj))
          console.log(typeof obj.metadata.labels)
          if ( "project" in this.form.metadata.labels){
            this.form.spec.contact = obj.metadata.labels["project"]
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
  components:{
    labelFrom
  },
  methods: {
    initContactGroups(){
      ListContactGroup().then(res=>{
        res.items.forEach(item =>{
          this.contactOptions.push(item.metadata.name)
        })
      })
      //console.log(this.contactOptions)
    },
    defaultSelect(){
      this.target_inpout_edit = this.form.spec.targets.length -1
    },
    delTarget(index){
      if ( this.form.spec.targets.length > 1 ){
        this.form.spec.targets.splice(index, 1)
      }
    },
    addHost(index){
      if (this.form.spec.targets[index] !== ""){
        this.form.spec.targets.push("")
        this.target_inpout_edit = index
      }
    },
    onSubmit() {
      this.$refs['from'].validate((valid) => {
        if (valid){
          this.form.metadata.labels["project"] = this.form.spec.contact
          if (this.isEdit){
            UpdateProbe(this.form).then(response => {
              this.onCancel()
              this.callBack()
            }).catch(err=>{
              console.log(err)
            })
          }else{
            console.log(this.form)
            CreateProbe(this.form).then(response => {
              this.onCancel()
              this.callBack()
            }).catch(err=>{
              console.log(err)
            })
          }
        }
      })
    },
    onCancel() {
      this.showDrawer = false
    },
    addField(items){
      items.push({key:"",value:""})
    },
    delField(index,items){
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

