<template>
  <el-form-item> 
        <el-col :span="5">
          <el-input v-model="data.key"/>
        </el-col>
        <el-col :span="2" class="line">-</el-col>
        <el-col :span="5">
          <el-input v-model="data.value"/>
        </el-col>
        <el-col :span="2" class="line">-</el-col>
        <el-col :span="5">
          asdfs
        </el-col>
</el-form-item>
</template>

<script>

export default {
  name: 'keyValue',
  props: {
      data: {
        type: Object,
        required: false
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
</style>
