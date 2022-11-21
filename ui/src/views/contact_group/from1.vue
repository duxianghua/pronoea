<template>
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
        <el-button type="primary" @click="onSubmit">Create</el-button>
        <el-button @click="onCancel" style="float:right;">Cancel</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>



<script>
import { CreateContactGroup } from '@/api/ContactGroup'
import { validEmail } from '@/utils/validate'
import { type } from 'os'
// import labelFrom from '@/components/labels'
export default {
  name: 'ContactGroupFrom1',
  data() {
    return {
      // tempLabels:[{key:"project",value:""}],
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
  created() {
        this.addItem(this.contactGroupItem.spec.members)
  },
  components:{
    },
  methods: {
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
          CreateContactGroup(this.contactGroupItem).then(resp=> {
            this.$store.commit('ChangeShowDrawer')
          }).catch(err=>{
            this.$message({
              message: err,
              type: "warning"
            })
          })
        }
      })
    },
    onCancel() {
      this.$store.commit('ChangeShowDrawer')
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

