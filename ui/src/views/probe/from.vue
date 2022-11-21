<template>
  <div class="app-container">
    <el-form ref="form" :model="form" label-width="120px">
      <el-form-item label="Name">
        <el-input v-model="form.metadata.name" />
      </el-form-item>
      <el-form-item label="Labels">
        <labelFrom :data="tempLabels">
        </labelFrom>
      </el-form-item>
      <el-form-item label="Targets" style="margin-top: 15px;">
          <el-container>
            <el-aside width="130px">
              <el-select v-model="form.spec.module.http.method">
                <el-option label="GET" value="GET"></el-option>
                <el-option label="POST" value="POST"></el-option>
              </el-select>
          </el-aside>
          <el-main style="padding: 0 0 0 15px;">
            <div style="margin: 0 0 0 0;" v-for="(item,index) in form.spec.targets" :key="index">
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
            </div>
          </el-main>
          </el-container>
          <el-collapse>
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
      <el-form-item>
        <el-button type="primary" @click="onSubmit">Create</el-button>
        <el-button @click="onCancel">Cancel</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>



<script>
import { CreateProbe } from '@/api/probe'
import labelFrom from '@/components/labels'
export default {
  name: 'probeForm',
  data() {
    return {
      tempLabels:[{key:"project",value:""}],
      headers: {},
      hosts: "",
      target_inpout_edit: -1,
      form: {
        "kind": "Probe",
        "apiVersion": "syncbug.io/v1",
        "metadata": {
            "name": "",
            "labels": {}
        },
      "spec": {
        "targets": [""],
        "labels": {
            "project": "baidu"
        },
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
    }
    }
  },
  computed:{
    // lablesData: {
    //   get(){
    //     this.tempLabels = Object.keys(this.form.metadata.labels).map((key)=>(
    //       {key: key, value: this.form.metadata.labels[key]}
    //     ))
    //     return this.tempLabels
    //   }
    // }
  },
  watch: {
      tempLabels: {
        handler(val){
        var _tempData = {}
        this.tempLabels.map(function(item){
            if (item.key.length > 2 && item.value.length > 2 ) {
              _tempData[item.key] = item.value
            }
         })
         this.form.metadata.labels = _tempData
     },
     deep: true
    },
    form: {
        handler(){
          this.target_inpout_edit = this.form.spec.targets.length -1
     },
     deep: true
    }
  },
  created() {
        this.defaultSelect()

  },
  components:{
    labelFrom
    },
  methods: {
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
    // initLables(){
    //   this.tempLabels = Object.keys(this.form.spec.labels).map((key)=>(
    //       {key: key, value: this.form.spec.labels[key]}
    //     ))
    // },
    onSubmit() {
      var ok = true
      if (this.form.spec.targets[this.form.spec.targets.length -1].length < 1){
        delete this.form.spec.targets[this.form.spec.targets.length -1]
      }
      this.form.spec.targets.forEach(element => {
        if (element === ""){
          this.$message({
            message: 'host is empty',
            type: 'warning'
          })
          ok=false
        }
      });
      if (ok){
        CreateProbe(this.form).then(response => {
          this.$store.commit('ChangeShowDrawer')
        }).catch(err=>{
          console.log(err)
        })
      }
    },
    onCancel() {
      this.$store.commit('ChangeShowDrawer')
    },
    addField(items){
      items.push({key:"",value:""})
    },
    delField(index,items){
      items.splice(index, 1)
    },
  }
}
</script>
@import url("//unpkg.com/element-ui@2.15.10/lib/theme-chalk/index.css");
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

