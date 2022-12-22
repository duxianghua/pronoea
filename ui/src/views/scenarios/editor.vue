<template>
  <div class="app-container">
    <!-- <el-header>
      <el-row :gutter="0" style="height: 45px;margin-bottom: 0px;">
        <el-col :span="18" >
          <el-form :inline="true" ref="form" :model="obj" :rules="formRules" hide-required-asterisk inline-message	>
            <el-form-item label="Name" prop="metadata.name">
              <el-input v-model="obj.metadata.name"/>
            </el-form-item>
            <el-form-item label="Interval" prop="metadata.labels.interval">
              <el-input v-model="obj.metadata.labels.interval"/>
            </el-form-item>
            <el-form-item label="ContactGroup" prop="metadata.labels.project">
            <el-select v-model="obj.metadata.labels.project" filterable placeholder="Please select">
              <el-option
                v-for="item in contactOptions"
                :key="item"
                :label="item"
                :value="item"
              />
              </el-select>
            </el-form-item>
          </el-form>
        </el-col>
        <el-col :span="3" >
            <el-button  style="width: 95%;" aria-disabled="true">SAVE</el-button>
        </el-col>
        <el-col :span="3" >
            <el-button  style="width: 95%;" >RUN</el-button>
        </el-col>
      </el-row>
    </el-header> -->
    <el-row>
      <el-col :span="20">
        <div ref="editor" class="editor" :style="{height: bodyHeight}">
          <div class="editor-header">
            <el-row :gutter="0">
              <el-col :span="18"><div class="grid-content">test.js</div></el-col>
              <el-col :span="6" >
                  <!-- <div >COPY SCRIPT</div> -->
                  <el-button  type="text" style="float: right;" >COPY SCRIPT</el-button>
              </el-col>
              <!-- <el-col :span="6"><div class="grid-content editor-header-button">COPY SCRIPT</div></el-col> -->
            </el-row>
          </div>
          <div class="editor-body">
            <MonacoEditor
              width="100%"
              :height=monacoEditorHeight
              theme="vs"
              v-model="obj.data['test.js']"
              @change="onChange"
              v-loading="loading"
            />
            <!-- <MonacoEditor
              width="100%"
              :height=monacoEditorHeight
              theme="vs"
              :value="obj.data['test.js']"
              v-model="obj.data['test.js']"
              @change="onChange"
              v-loading="loading"
            /> -->
          </div>
        </div>
      </el-col>
      <el-col :span="4">
        <div class="editor-sidebar">
          <el-row>
            <el-button :disabled="saveDisabled" @click="onSubmit">SAVE</el-button>
          </el-row>
          <el-divider></el-divider>
          <el-form ref="form"
            :model="obj" 
            :rules="formRules"
            label-position="top" 
            hide-required-asterisk
          >
            <el-form-item label="Name" prop="metadata.name" >
              <el-input v-model="obj.metadata.name" :disabled="this.viewInfo.name?true:false"/>
            </el-form-item>
            <el-form-item label="Interval" prop="data.interval" name="asdf">
              <el-input v-model="obj.data.interval" @change="onChange"/>
            </el-form-item>
            <el-form-item label="ContactGroup" prop="metadata.labels.project">
            <el-select v-model="obj.metadata.labels.project" filterable placeholder="Please select"  @change="onChange">
              <el-option
                v-for="item in contactOptions"
                :key="item"
                :label="item"
                :value="item"
              />
              </el-select>
            </el-form-item>
          </el-form>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import MonacoEditor from '@/components/MonacoEditor'
import { GetScenarios, CreateScenarios, statusScenarios, UpdateScenarios } from '@/api/scenarios'
import { ListContactGroup } from '@/api/ContactGroup'

export default {
  components: {
    MonacoEditor
  },
  data() {
    return {
      formRules: {
        'metadata.name': [
          { required: true, message: 'please enter name', trigger: 'blur' },
          { min: 3, message: 'min size is 3', trigger: 'blur' }
        ],
        'metadata.labels.project': [
          { required: true, message: 'please select project', trigger: 'change' }
        ],
        'data.interval': [
          { name: "asfdd", required: true, type: "integer", min: 10, trigger: ['blur', 'change'], transform(value) {return Number(value)} },
        ]
      },
      saveDisabled: true,
      contactOptions: [],
      obj: {
        metadata: {
          name: "",
          namespace: "",
          labels: {
            "app.kubernetes.io/component": "scenarios",
            "app.kubernetes.io/managed-by": "pronoea",
            "project": "",
            "interval": "60",
          }
        },
        data: {
          "test.js": ""
        },
      },

      code: "",
      loading: true,
      // style: {
      //   width: "100%",
      //   height: "600"
      // },
      monacoEditorHeight: 300,
      bodyHeight: 300,
      codeTemp: `import { check } from 'k6';
import http from 'k6/http';

export default function () {
  const res = http.get('http://test.k6.io/');
  check(res, {
    'is status 200': (r) => r.status === 200,
  });
}`,
      viewInfo: {
        namespace: "",
        name: ""
      }
    }
  },
  created() {
    this.bodyHeight = (window.innerHeight - 100) + 'px'
    this.monacoEditorHeight = window.innerHeight - 155
    this.viewInfo.namespace = this.$route.params.namespace
    this.viewInfo.name = this.$route.params.name
    this.fetchData()
    this.initContactGroups()
  },
  mounted(){
    window.addEventListener('resize', this.getRefs);
  },
  methods: {
    getRefs(){
      this.bodyHeight = (window.innerHeight - 100) + 'px'
      this.monacoEditorHeight = this.$refs.editor.offsetHeight - 52
    },
    onChange(value) {
      if (value !== this.obj.data['test.js']){
        this.saveDisabled = false
      }else{
        this.saveDisabled = true
      }
    },
    fetchData(){
      if (this.viewInfo.namespace && this.viewInfo.name){
        GetScenarios(this.viewInfo.namespace, this.viewInfo.name).then(response => {
          this.obj = response
          this.code = String(response.data["test.js"])
          this.loading = false
          
        }).catch(err => {
          console.log(err)
          this.loading = false
        })
      }else{
        this.obj.data['test.js'] = this.codeTemp
        this.loading = false
      }
    },
    onSubmit() {
      this.$refs['form'].validate((valid) => {
        console.log(this.obj)
        if (valid) {
          if (this.obj.metadata.resourceVersion) {
            UpdateScenarios(this.obj).then(response => {
              this.saveDisabled = true
            }).catch(err => {
              console.log(err)
            })
          } else {
            CreateScenarios(this.obj).then(response => {
              this.saveDisabled = true
            }).catch(err => {
              console.log(err)
            })
          }
        }
      })
    },
    initContactGroups() {
      ListContactGroup().then(res => {
        res.items.forEach(item => {
          this.contactOptions.push(item.metadata.name)
        })
      })
    },
  },
}
</script>

<style scoped>
.app-container{
  height: 100%;
}
.editor {
  /* max-width: 80%; */
  width: 100%;
  /* min-height: 750px; */
  height: 100%;
  /* box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04) */
  box-shadow: 0px 0px 10px #dedede;
  /* position: absolute; */
  display:inline-block
}
.editor-header{
  min-height: 38px;
  /* box-shadow: 0 2px 4px rgba(0, 0, 0, .12), 0 0 6px rgba(0, 0, 0, .04); */
  /* display: flex; */
  -moz-box-align: center;
  padding: 5px 5px 5px 10px;
  color: rgb(77, 79, 92);
  border-bottom: 2px solid rgb(121, 121, 121);
}
.editor-body{
  min-height: 350px;
}
.editor-sidebar{
  width: 100%;
  text-align: center;
  padding: 0px 20px 20px 20px;
  border-radius: 0px;
}
.el-button--default{
  min-width: 100px;
  margin-bottom: 0px;
  width: 100%;
}
.editor-input {
  border-radius: 0px !important;
    /* -webkit-appearance: none;
    background-color: #FFF;
    background-image: none;
    
    border: 1px solid #DCDFE6;
    -webkit-box-sizing: border-box;
    box-sizing: border-box;
    color: #606266;
    display: inline-block;
    font-size: inherit;
    height: 40px;
    line-height: 40px;
    outline: 0;
    padding: 0 15px;
    -webkit-transition: border-color .2s cubic-bezier(.645,.045,.355,1);
    transition: border-color .2s cubic-bezier(.645,.045,.355,1);
    width: 10%; */
}
.editor-header-button{
  position: absolute;
  width: 100%;
}
.el-form-item__content{
  margin-left: 0px;
}
.element.style {
  margin-left: 0px;
}

.grid-content{
  padding: 8px 14px 8px 14px;
  color: rgb(0, 0, 0);
  font-weight: 600;
  line-height: 1.75;
  font-family: TTNormsPro;
  min-height: 36px;
  font-size: 12px;
}
.el-button{
  color: rgb(125, 100, 255);
  border-radius: 0px;
  border-color:rgb(125, 100, 255); border-width: 2px;
}
.el-button.is-disabled{
  color: #C0C4CC;
  border-color:#EBEEF5; border-width: 2px;
}

::v-deep .el-input__inner {
  /* color: rgb(125, 100, 255); */
  border-radius: 0px !important;
  /* border-color:rgb(125, 100, 255); border-width: 1px; */
}
::v-deep .el-form-item__label{
  text-align: right;
  float: left;
  font-size: 14px;
  padding: 0 12px 0 0;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
}
.el-button--text{
  padding: 10px 10px 10px 10px;
  color: rgb(109, 52, 241);
  border-color:rgb(125, 100, 255); border-width: 0px;
}

.el-select {
  width: 100%;
}
</style>
