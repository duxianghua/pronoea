<template>
  <div>
    <el-row v-for="(item,index) in value" :key="index">
      <el-col :span="10">
        <el-input v-model="value[index].key" />
      </el-col>
      <el-col :span="1" class="line">
        <i />
      </el-col>
      <el-col :span="10">
        <el-input v-model="value[index].value" style="margin-left: 10px;" />
      </el-col>
      <el-col :span="3">
        <svg v-show="value.length - 1 == index" viewBox="0 0 24 24" width="32" height="32" @click="addField(value)">
          <path fill="none" d="M0 0h24v24H0z" />
          <path
            fill="green"
            d="M11 11V7h2v4h4v2h-4v4h-2v-4H7v-2h4zm1 11C6.477 22 2 17.523 2 12S6.477 2 12 2s10 4.477 10 10-4.477 10-10 10zm0-2a8 8 0 1 0 0-16 8 8 0 0 0 0 16z"
          />
        </svg>
        <svg viewBox="0 0 24 24" width="32" height="32" @click="delField(index, value)">
          <path fill="none" d="M0 0h24v24H0z" />
          <path
            fill="red"
            d="M12 22C6.477 22 2 17.523 2 12S6.477 2 12 2s10 4.477 10 10-4.477 10-10 10zm0-2a8 8 0 1 0 0-16 8 8 0 0 0 0 16zm0-9.414l2.828-2.829 1.415 1.415L13.414 12l2.829 2.828-1.415 1.415L12 13.414l-2.828 2.829-1.415-1.415L10.586 12 7.757 9.172l1.415-1.415L12 10.586z"
          />
        </svg>
      </el-col>
    </el-row>
  </div>

</template>

<script>

export default {
  name: 'LabelFrom',
  props: {
    data: {
      type: Array,
      required: true
    }
  },
  computed: {
    value: {
      get() {
        return this.data
      }
    },
    addFieldShow(index) {
      console.log(index)
      this.data.length === index
    }
  },
  methods: {
    formatData() {
      console.log(data)
      this.data
    },
    onSubmit() {
      this.$message('asf')
      console.log(this.form)
    },
    onCancel() {
      this.$message({
        message: 'cancel!',
        type: 'warning'
      })
    },
    hander() {
      this.$message('asf')
      console.log(this.form)
    },
    addField(items) {
      if (items[items.length - 1].key !== '') {
        items.push({ key: '', value: '' })
      }
    },
    delField(index, items) {
      if (items.length === 1 && items[0].key !== '' && items[0].value !== '') {
        items.splice(index, 1)
        items.push({ key: '', value: '' })
      }
      if (items.length > 1) {
        items.splice(index, 1)
      }
    }
  }

}
</script>

<style scoped>
.svg-icon {
  width: 1em;
  height: 1em;
  vertical-align: -0.15em;
  fill: currentColor;
  overflow: hidden;
}

.svg-external-icon {
  background-color: currentColor;
  mask-size: cover!important;
  display: inline-block;
}
.el-col{
  margin-bottom: 10px;
  text-align: right;
}
.svg-icon{
  text-align: right;
}
</style>
