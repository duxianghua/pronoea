<template>
    <el-drawer
      title="Add Contact Group"
      :visible.sync="showDrawer"
      v-if="showDrawer"
      direction="rtl"
      :modal="true"
      :show-close="true"
      size="30%"
      :before-close="onCancel"
      :destroy-on-close="true"
    >
  <div class="app-container">
    <el-form ref="from" :model="contactGroupItem" :rules="formRules" label-width="120px">
      <el-form-item label="Name" prop="metadata.name">
        <el-input v-model="contactGroupItem.metadata.name" />
      </el-form-item>
      <el-form-item label="Projects" prop="spec.projects">
        <el-select v-model="contactGroupItem.spec.projects" filterable placeholder="Please select">
          <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.label"
            :value="item.value">
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="Members" >
        <div v-for="(email,index) in contactGroupItem.spec.members" :key="index">
          <el-form-item :prop="'spec.members.'+index" :rules="formRules.email">
              <el-row >
                  <el-input
                    placeholder="devops@email.com"
                    v-model="contactGroupItem.spec.members[index]"
                    class="input-with-select" 
                    style="margin-bottom: 10px;"
                  >
                    <i slot="suffix" class="el-icon-delete" @click="delItem(contactGroupItem.spec.members, index)"></i>
                  </el-input>
                </el-row>
                
      </el-form-item>
        </div>

      <el-row style="float:right; margin-top: -10px">
                  <el-button type="text" @click="addItem(contactGroupItem.spec.members)">Add Member</el-button>
                </el-row>
    </el-form-item>
      <el-form-item style="margin-top: 100px;">
        <el-button type="primary" @click="onSubmit">Apply</el-button>
        <el-button @click="onCancel" style="float:right;">Cancel</el-button>
      </el-form-item>
    </el-form>
  </div>
</el-drawer>
</template>



<script>
import { CreateContactGroup, UpdateContactGroup } from '@/api/ContactGroup'

// import labelFrom from '@/components/labels'
export default {
  name: 'ContactGroupFrom',
  props: {
    isActive: {
      type: Boolean,
      default: false
    },
    isEdit: {
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
    return {
      showDrawer: this.isActive,
      formRules: {
        'metadata.name': [
          {required: true, message: "please enter name", trigger: "blur"},
          {min: 3, message: "min size is 3", trigger:'blur'}
          
        ],
        'spec.projects': [
          {required: true, message: "please select project", trigger: "change"},          
        ],
        'email': [
        {type: 'email', message: "please enter email address", trigger: ["blur", 'change']}
        ]
      },
      headers: {},
      hosts: "",
      data: [{"key":"devops", "label": "devops"}, {"key":"sso", "label": "sso"}],
      target_inpout_edit: -1,
      contactGroupItem: {
        "kind": "ContactGroup",
        "apiVersion": "pronoea.io/v1",
        "metadata": {
            "name": "",
            "labels": {}
        },
        "spec": {
          "members": [],
          "projects": ""
        },
      },
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
  watch:{
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
    formData: {
      handler: function(val, oldval){
        this.contactGroupItem = JSON.parse(JSON.stringify(val))
      },
      immediate: true
    }
    // isEdit: {
    //   handler: function(val, oldval){
    //     this.contactGroupItem = this.data
    //   },
    //   immediate: true
    // }
  },
  // computed: {
  //   contactGroupItem(){
  //     console.log("computed:contactGroupItem")
  //     return JSON.parse(JSON.stringify(this.formData))
  //   }
  // },
  created() {
        //this.testFunc()
        //this.addItem(this.contactGroupItem.spec.members)
  },
  components:{
    },
  methods: {
    testFunc(){
      console.log("testFunc")
      console.log(this.formData)
      this.contactGroupItem = this.data
    },
    delItem(data, index){
      if ( data.length > 1 ){
        delete data[index]
      }
    },
    addItem(data){
      data.push("")
    },
    onSubmit() {
      this.$refs['from'].validate((valid) => {
        if (valid){
          if (this.isEdit){
            UpdateContactGroup(this.contactGroupItem).then(resp=> {
              this.onCancel()
              this.callBack()
          }).catch(err=>{
              this.$message({
                message: err,
                type: "warning"
              })
          })
          }else{
              CreateContactGroup(this.contactGroupItem).then(resp=> {
                this.onCancel()
                this.callBack()
              }).catch(err=>{
                this.$message({
                  message: err,
                  type: "warning"
                })
              })
          }

        }
      })
    },
    onCancel() {
      this.showDrawer = false
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

