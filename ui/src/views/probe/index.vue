<template>
  <div class="app-container">
    <el-row style="margin-bottom: 20px;">
      <el-button @click="DrawerSwitch">
        Add Probe
    </el-button>
    </el-row>
    <!-- Table -->
    <el-table
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
    >
      <el-table-column align="center" label="ID" width="95">
        <template slot-scope="scope">
          {{ scope.$index }}
        </template>
      </el-table-column>
      <el-table-column label="Name" width="260">
        <template slot-scope="scope">
          {{ scope.row.metadata.name }}
        </template>
      </el-table-column>
      <el-table-column label="Targets" align="center">
        <template slot-scope="scope">
          <p v-for="v in scope.row.spec.targets">
            {{ v }}
          </p>
        </template>
      </el-table-column>
      <el-table-column label="ProbeType" width="110" align="center">
        <template slot-scope="scope">
          {{ scope.row.spec.module.prober }}
        </template>
      </el-table-column>
      <el-table-column label="Monitoring" width="110" align="center">
        <template slot-scope="scope">
          <el-switch
            v-model="scope.row.spec.pause "
            active-color="#13ce66"
            inactive-color="#ff4949"
            @change="pauseSwitch(scope.row, scope.$index)">
          </el-switch>
        </template>
      </el-table-column>
      <el-table-column align="center" prop="created_at" label="CreateTime" width="200">
        <template slot-scope="scope">
          <!-- <i class="el-icon-time"/> -->
          <span>{{ scope.row.metadata.creationTimestamp }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" prop="created_at" label="Aciton" width="200">
        <template slot-scope="scope">
          <el-button type="primary" @click="probeStatus(scope.row.metadata)">Check</el-button>
          <el-button type="danger" icon="el-icon-delete" circle size="mini" @click="delProbe(scope.row.metadata)"></el-button>
        </template>
        
      </el-table-column>
    </el-table>
    <!-- FORM -->
    <el-drawer
      title="Add HTTP Probe"
      :visible.sync="showDrawer"
      v-if="showDrawer"
      direction="rtl"
      :modal="true"
      :show-close="true"
      :wrapperClosable="false"
      size="60%"
      :before-close="closeHander"
      :destroy-on-close="true"
    >
      <probeForm></probeForm>
    </el-drawer>

    <el-dialog
    title="提示"
    :visible.sync="centerDialogVisible"
    width="80%"
    center>
    <template>
      <el-tabs >
          <el-tab-pane  v-for="valuea, key, index in statusContext" :label="key" :key="index" >
          <el-input
          type="textarea"
          :rows="18"
          :value=valuea>
          </el-input>
        </el-tab-pane>
      </el-tabs>
    </template>
      <!-- <el-row v-for="valuea, key, index in statusContext">
        <span>{{ key }}</span>
        
      </el-row> -->
      <span slot="footer" class="dialog-footer">
      <el-button @click="centerDialogVisible = false">取 消</el-button>
      <el-button type="primary" @click="centerDialogVisible = false">确 定</el-button>
    </span>
  </el-dialog>
  </div>
</template>

<script>
import { getList, deleteProbe, statusProbe, UpdateProbe } from '@/api/probe'
import { get } from 'js-cookie'
import probeForm from './from'

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
  data() {
    return {
      centerDialogVisible: false,
      statusContext: {},
      list: null,
      listLoading: true
    }
  },
  created() {
    this.fetchData()
  },
  computed:{
    showDrawer(){
        return this.$store.state.probe.showDrawer
    }
  },
  components:{
    probeForm
  },
  methods: {
    fetchData() {
      this.listLoading = true
      getList().then(response => {
        console.log(response)
        this.list = response.items
        this.listLoading = false
      }).catch(err=>{
        console.log(err)
      })
    },
    delProbe(metadata){
      console.log(metadata)
      deleteProbe(metadata.namespace, metadata.name).then(response => {
        this.fetchData()
      }).catch(err=>{
        console.log(err)
      })
    },
    probeStatus(metadata){
      this.centerDialogVisible=true
      console.log(metadata)
      statusProbe(metadata.name, {namespace: metadata.namespace}).then(response => {
        this.statusContext = response
      }).catch(err=>{
        console.log(err)
      })
    },
    DrawerSwitch() {
      this.$store.commit('ChangeShowDrawer')
    },
    closeHander(){
      this.DrawerSwitch()
      this.fetchData()
    },
    pauseSwitch(obj, index){
      UpdateProbe(obj).then(response=>{
        console.log(index)
        console.log(response)
        this.list[index] = response
      }).catch(err=>{
        this.$message(err.$message)
      })
    }
  }
}
</script>
