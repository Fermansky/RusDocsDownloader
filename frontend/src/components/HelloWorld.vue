<script setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {GetPdf, OpenFileOrFolder, Scrape} from '../../wailsjs/go/main/App'
import {EventsOff, EventsOn, EventsOnce} from "../../wailsjs/runtime/runtime.js";
import {ElMessage, ElMessageBox} from "element-plus";

const startPage = ref(0)
const endPage = ref(0)
const fixedStartPage = ref(0)
const fixedEndPage = ref(0)
const bookName = ref("")
const quality = ref(8)
const processing = ref(false)

const progress = ref(0)
const tip = ref("")
const url = ref("")
const mode = ref(0)
const downloading = ref(false)

const bookInfo = reactive({
  title: "Политбюро и дело Виктора Абакумова",
  pageNum: 752,
  pages: []
})

const getPagePreview = (id, quality) => {
  if (!Number.isInteger(id)) {
    return ''
  }
  return `https://docs.historyrussia.org/pages/${id}/zooms/${quality}`
}
const onPagePreviewLoaded = () => {

}

const sanitizeFileName = (name) => {
  return name.replace(/[<>:"\/\\|?*]+/g, '_')
}

function pageChanged() {
  const start = parseInt(startPage.value)
  const end = parseInt(endPage.value)

  fixedStartPage.value = start
  fixedEndPage.value = end

  if (end < start) {
    fixedEndPage.value = start
    endPage.value = fixedEndPage.value
  }
}

const downloadCheck = () => {
  if (Math.abs(bookInfo.pageNum-bookInfo.pages.length) > 50) {
    ElMessageBox.confirm(
        '书籍页数和下载页数相差较大，确定下载吗？',
        '注意',
        {
          confirmButtonText: '下载',
          cancelButtonText: '取消',
          type: 'warning',
        }
    )
        .then(() => {
          download()
        })
        .catch(() => {
        })
  } else {
    download()
  }
}

const download = async () => {
  EventsOnce("complete", path => {
    processing.value = false
    downloading.value = false
    EventsOff("progress")
    tip.value = path
  })
  EventsOn("progress", p => {
    progress.value = parseInt(p)
  })
  const sanitizedBookName = sanitizeFileName(bookInfo.title)
  processing.value = true
  downloading.value = true
  tip.value = ""
  await GetPdf(bookInfo.pages, sanitizedBookName, quality.value)
  progress.value = 0
  ElMessage.success("下载完成")
}

const openFile = file => {
  OpenFileOrFolder(file)
}

const scrape = async () => {
  processing.value = true
  Scrape(url.value).then(book => {
    if (book.Type === 1) {
      startPage.value = book.Pages[0]
      endPage.value = book.Pages[book.Pages.length-1]
      bookName.value = book.Name
      pageChanged()
    }
    bookInfo.title = book.Name
    bookInfo.pages = book.Pages
    bookInfo.pageNum = book.PageNum
    console.log(bookInfo.pages)
    thumbnail.value = getPagePreview(bookInfo.pages[0], 5)
    processing.value = false
  })
}

const thumbnail = ref("")

</script>

<template>
  <main>
<!--    <div style="margin: 20px 0 ">-->
<!--      <el-radio-group v-model="mode" size="large">-->
<!--        <el-radio-button label="解析模式" :value="0" />-->
<!--        <el-radio-button label="手动模式" :value="1" />-->
<!--      </el-radio-group>-->
<!--    </div>-->
    <div v-if="mode === 0" >
      <el-input v-model="url" style="margin: 10px 0; padding: 0 20px" :disabled="processing" placeholder="请输入网页地址">

      </el-input>
      <el-button size="large" type="primary" @click="scrape" :disabled="processing">
        解析
      </el-button>

      <div style="display: flex; padding: 20px; align-items: center; justify-content: center" v-if="thumbnail">
        <el-image :src="thumbnail"
         class="book-thumbnail"/>
        <div style="align-content: center; text-align: left; margin-left: 2%">
          <div style="font-weight: bold; font-size: 20px">书籍信息</div>
          <div>
            书籍名称：{{bookInfo.title}}
          </div>
          <div>
            书籍页数：{{bookInfo.pageNum}}
          </div>
          <div>
            下载页数：{{bookInfo.pages.length}}
          </div>
          <div style=" display: flex; align-items: center; margin-top: 10px">
            <div style="flex: none;">
              图片质量：
            </div>
            <el-slider :disabled="processing" v-model="quality" :step="1" :min="0" :max="8" show-stops show-input style="min-width: 300px; width: 50vw"/>
          </div>
          <el-button @click="downloadCheck" size="large" type="success" :disabled="processing">
            下载
          </el-button>
        </div>
      </div>
    </div>
    <div v-if="mode === 1">
      <div style="display: flex; justify-content: center" >
        <div style="max-width: 120px">
          <div style="height: 150px">
            <el-image :src=getPagePreview(fixedStartPage,0) class="preview-page" @load="onPagePreviewLoaded"/>
            <div>{{fixedStartPage}}</div>
          </div>
          <el-input :disabled="processing" v-model="startPage" type="number" min="0" step="1" :formatter="value => value.replace(/^0+(?=\d)/, '')"
                    @change="pageChanged" placeholder="请输入起始页码"/>
        </div>
        <div style="margin: 20px; min-width: 20px">
          <div v-if="fixedStartPage && fixedEndPage">{{fixedEndPage-fixedStartPage+1}}页</div>
        </div>
        <div style="max-width: 120px">
          <div style="height: 150px">
            <el-image :src=getPagePreview(fixedEndPage,0) class="preview-page" @load="onPagePreviewLoaded"/>
            <div>{{fixedEndPage}}</div>
          </div>
          <el-input :disabled="processing" v-model="endPage" type="number" :formatter="value => value.replace(/^0+(?=\d)/, '')"
                    :min="fixedStartPage" step="1" @change="pageChanged" placeholder="请输入结束页码"/>
        </div>
      </div>
      <div style="margin: 20px;">
        <el-input :disabled="processing" v-model="bookName" placeholder="请输入书名"></el-input>
      </div>
      <div style="padding: 0 20px; display: flex; align-items: center">
        <div style="flex: none;">
          图片质量：
        </div>
        <el-slider :disabled="processing" v-model="quality" :step="1" :min="0" :max="8" show-stops show-input />
      </div>
      <div>
        <el-button :disabled="processing" type="success" size="large" @click="download">下载</el-button>
      </div>
    </div>

    <div style="padding: 0 20px; display: flex" v-if="downloading">
      <div>
        下载进度：
      </div>
      <el-progress
          style="flex: auto"
          :percentage="progress"
          :stroke-width="15"
          striped
          striped-flow
          :duration="10"
      />
    </div>

    <div>
      <div @click="openFile(tip)" class="file-link">
        {{tip}}
      </div>
    </div>
  </main>
</template>

<style scoped>

.descriptions :deep(.el-descriptions__title) {
  color: white;
  font-size: 16px;
  font-weight: bold;
}

.file-link {
  padding: 10px;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.file-link:hover {
  color: #67c23a;
  cursor: pointer;
}

.preview-page {
  height: 120px;
  min-height: 120px;
  max-width: 120px;
  min-width: 70px;
}

.book-thumbnail {
  height: 240px;
  min-height: 240px;
  max-width: 240px;
  min-width: 50px;
}

</style>
