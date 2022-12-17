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
      border
      fit
      size="nini"
      highlight-current-row
      style="min-height: 540px; box-shadow: 0px 0px 10px #dedede;"
    >
      >
      <el-table-column align="center" label="ID" width="95">
        <template slot-scope="scope">
          {{ scope.$index }}
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
          {{ scope.row.data.interval }}s
        </template>
      </el-table-column>
      <el-table-column label="ProbeType" width="110" align="center">
        <template>
          k6
        </template>
      </el-table-column>
      <el-table-column align="center" prop="created_at" label="CreateTime" width="200">
        <template slot-scope="scope">
          <!-- <i class="el-icon-time"/> -->
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
import { ListScenarios, DeleteScenarios, statusScenarios, UpdateScenarios } from '@/api/scenarios'

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
      formEdit: false,
      formData: {},
      showDrawer: false,
      centerDialogVisible: false,
      statusContext: {},
      list: null,
      listKey: Math.random(),
      listLoading: true,
      loading: true
    }
  },
  computed: {
    // showDrawer(){
    //     return this.$store.state.probe.showDrawer
    // }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      ListScenarios().then(response => {
        this.list = response.items
        this.listLoading = false
      }).catch(err => {
        console.log(err)
      })
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
    pauseSwitch(scope) {
      UpdateScenarios(scope.row).then(response => {
        this.$set(this.list, scope.$index, response)
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
  }
}
</script>
