<template>
  <div class="app-container">
    <div style="margin-bottom: 5px; box-shadow: 0px 0px 10px #dedede; height: 60px; padding:10px">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-input
            placeholder="Search"
            prefix-icon="el-icon-search"
          />
        </el-col>
        <el-col :span="8"><div class="grid-content bg-purple" /></el-col>
        <el-col :span="8" style="float:right">
          <el-button type="primary" style="float:right" @click="gotoEditor({})">
            Create Scenarios
          </el-button>
        </el-col>

      </el-row>
    </div>
    <el-table
      :key="listKey"
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      size="medium"
      style="min-height: 450px; box-shadow: 0px 0px 10px #dedede;"
      :row-class-name="callRowClass"

    >
      >
      <el-table-column width="15" type="expand">
        <template slot-scope="scope">
          <el-table
            :data="getChecksRate(scope.row)"
            element-loading-text="Loading"
            v-loading="checksLoading"
            :cell-class-name="checksClass"
            class="test123"
          >
          <el-table-column width="3"></el-table-column>
            <el-table-column label="Name">
              <template slot-scope="scope">
                {{ scope.row.metric.check }}
              </template>
            </el-table-column>
            <el-table-column label="Group">
              <template slot-scope="scope">
                {{ scope.row.metric.group }}
              </template>
            </el-table-column>
            <el-table-column label="State">
              <template slot-scope="scope">
                <i class="el-icon-close" style="color: red;" v-if="scope.row.value[1]==='0'"></i>
                <i class="el-icon-check" style="color: green;" v-if="scope.row.value[1]==='1'"></i>
              </template>
            </el-table-column>
        </el-table>
      </template>
      </el-table-column>
      <el-table-column label="Name" width="260">
        <template slot-scope="scope">
          <el-button type="text" @click="gotoEditor(scope.row)">
            {{ scope.row.metadata.name }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column label="Interval" align="center">
        <template slot-scope="scope">
          {{ scope.row.status.interval }}
        </template>
      </el-table-column>
      <el-table-column label="ProbeType"  align="center">
        <template>
          k6
        </template>
      </el-table-column>
      <el-table-column label="Checks" width="110" align="center">
        <template slot-scope="scope">
          <el-tag class="tagClass" v-if="scope.row.status.health === 'success'" type="success" effect="dark" >
            {{ scope.row.status.health }}
          </el-tag>
          <el-tag class="tagClass" v-if="scope.row.status.health === 'unknown'" type="info" effect="dark" >
            {{ scope.row.status.health }}
          </el-tag>
          <el-tag class="tagClass" v-if="scope.row.status.health === 'fail'" type="danger" effect="dark" >
            {{ scope.row.status.health }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Monitoring" width="110" align="center">
        <template slot-scope="scope">
          <el-switch
            v-model="scope.row.status.paused"
            inactive-value="true"
            active-value="false"
            
            inactive-color="#ff4949"
            active-color="#13ce66"
            @change=pausedSwitch(scope)
          />
        </template>
      </el-table-column>
      <el-table-column align="center" prop="created_at" label="CreateTime" width="200">
        <template slot-scope="scope">
          <span>{{ scope.row.metadata.creationTimestamp }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" prop="created_at" label="Aciton" width="210">
        <template slot-scope="scope">
          <!-- <el-button type="primary" size="mini" @click="probeStatus(scope.row)">Check</el-button> -->
          <el-button type="danger" icon="el-icon-delete" size="mini" style="width: 60%;" @click="deleteScenarios(scope)" />
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { ListScenarios, DeleteScenarios, StatusScenarios, PatchScenarios } from '@/api/scenarios'
import { nullLiteral } from '@babel/types'

export default {
  filters: {
    statusFilter(status) {
      const statusMap = {
        published: 'success',
        draft: 'gray',
        deleted: 'danger'
      }
      return statusMap[status]
    }
  },
  components: {
  },
  data() {
    return {
      checksList: [
        {"check": "asdf", "value": ""},
      ],
      checksLoading: true,
      formEdit: false,
      formData: {},
      showDrawer: true,
      centerDialogVisible: false,
      statusContext: {},
      list: null,
      listKey: Math.random(),
      listLoading: true,
      loading: true
    }
  },
  computed: {
  },
  created() {
    this.fetchData()
  },
  methods: {
    callRowClass(obj){
      if (obj.row.status.health === "success"){
          return 'green'
        }
        if (obj.row.status.health === "fail"){
          return 'red'
        }
        if (obj.row.status.health === "unknown"){
          return 'unknown'
        }
    },  
    testClass(obj){
      if (obj.columnIndex == 0){
        if (obj.row.status.health === "success"){
          return 'greenStatus'
        }
        if (obj.row.status.health === "fail"){
          return 'redStatus'
        }
        if (obj.row.status.health === "unknown"){
          return 'unknownStatus'
        }
      }
      
    },
    checksClass(obj){
      if (obj.columnIndex == 0){
        if (obj.row.value[1] === "1"){
          return 'greenStatus'
        }
        if (obj.row.value[1] === "0"){
          return 'redStatus'
        }
      }
      
    },
    fetchData() {
      this.listLoading = true
      ListScenarios().then(response => {
        this.list = response.items
        this.listLoading = false
      }).catch(err => {
        this.$message(err.$message)
      })
    },
    getChecksRate(row) {
      if (!row.checks){
        row.checks = []
        this.checksLoading=true
          StatusScenarios(row.metadata.name, {namespace: row.metadata.namespace}).then(response=>{
            row.checks = response
            this.checksLoading=false
          }).catch(err => {
            this.$message(err.$message)
          })
        }
        return row.checks
    },
    deleteScenarios(scope) {
      this.$confirm('Are you sure to delete permanently', {
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        DeleteScenarios(scope.row.metadata.name, { namespace: scope.row.metadata.namespace }).then(response => {
          this.list.splice(scope.$index, 1)
          this.$message({
            type: 'success',
            message: '删除成功!'
          })
        }).catch(err => {
          console.log(err)
        })
      })
    },
    probeStatus(row) {
      this.statusContext = {}
      this.centerDialogVisible = true
      this.loading = true
      row.spec.targets.forEach(element => {
        this.statusContext[element] = ''
      })
      this.loading = true
      statusProbe(row.metadata.name, { namespace: row.metadata.namespace }).then(response => {
        this.statusContext = response
        this.loading = false
      }).catch(err => {
        console.log(err)
        this.loading = false
      })
    },
    DrawerSwitch() {
      this.$store.commit('ChangeShowDrawer')
    },
    closeHander() {
      this.DrawerSwitch()
      this.fetchData()
    },
    pausedSwitch(scope) {
      PatchScenarios(scope.row.metadata.name, {"paused": scope.row.status.paused, namespace: scope.row.metadata.namespace}).then(response => {
        //this.$set(this.list, scope.$index, response)
      }).catch(err => {
        this.$message(err.$message)
      })
    },
    gotoEditor(data){
      if (Object.keys(data).length === 0 ){
        this.$router.push({
          path: `scenarios/editor`
        })
      }else{
        this.$router.push({
          path: `scenarios/editor/${data.metadata.namespace}/${data.metadata.name}`
        })
      }
    },

    expandChange(data){
      console.log(data)
      //this.getChecksRate(data)
    },
  }
}
</script>

<style>
/* .el-table__expanded-cell[class*=cell] {
  padding: 5px 60px 20px 60px;
  min-width: 0;
  box-sizing: border-box;
  text-overflow: ellipsis;
  vertical-align: middle;
  position: relative;
  text-align: left;
} */

.testClass {
  background-color: green;
}

.greenStatus {
  background-color: rgb(134, 188, 83);
}

.redStatus {
  background-color: red;
}

.unknownStatus {
  background-color: darkgrey;
}

.el-table td, .el-table th{
  padding: 6px 0;
}

.green::after{
  content: "";
  width: 4px;
  top: 0px;
  left: 0px;
  height: 100%;
  transition: all 120ms ease 0s;
  position: absolute;
  background-color: rgb(134, 188, 83);
}

.red::after{
  content: "";
  width: 4px;
  top: 0px;
  left: 0px;
  height: 100%;
  transition: all 120ms ease 0s;
  position: absolute;
  background-color: rgb(254, 59, 59);
}
.unknown::after{
  content: "";
  width: 4px;
  top: 0px;
  left: 0px;
  height: 100%;

  position: absolute;
  background-color:#909399;
}

.green{
  opacity: 0.8;
  padding-left: 5px;
  position: relative;
}

.unknown{
  opacity: 0.8;
  padding-left: 5px;
  position: relative;
}

.red{
  opacity: 0.8;
  padding-left: 5px;
  position: relative;
}

.tagClass{
  width: 80px;
  font-size: 12px;
  font-weight:bold;
}
</style>