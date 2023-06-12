<template>
	<div class="list">
		<div>
			<ul :style="`transform: translateY(${clientHeightTotal}px);`">
				<li v-for="(item, key) in dataList" :key="key" :id="`list${key}`" style="height: 100px; overflow: hidden;">
					{{ item }}</li>
			</ul>
		</div>
	</div>
</template>

<script setup>
import { nextTick, onMounted, onUpdated, ref, watch } from 'vue';

// 列表数据
const dataList = ref([])
// 列表 预览的数据
const preListBefore = ref([])
const preListAfter = ref([])

const total = ref(0)
const curPage = ref(1)



const clientHeightTotal = ref(0)


// 创建观察者
let b = false
onUpdated(async () => {

	if (b) {
		return
	}

	// 创建观察者
	const observer = new IntersectionObserver((changes) => {
		// 元素被隐藏
		if (!changes[0].isIntersecting) {
			// console.log("有一个元素被隐藏了 ");


			updateListAfter()
		}
	})
	observer.observe(document.querySelector("#list1"))


	// 创建观察者二
	const observer2 = new IntersectionObserver(async (changes) => {
		// 元素被隐藏
		if (changes[0].isIntersecting) {
			// console.log("元素被显示了");

			updateListBefore()
		}
	})

	observer2.observe(document.querySelector("#list0"))

	b = !b
})



// 更新列表
function updateListAfter() {
	const first = document.querySelector("#list0")

	// 计算元素移动的距离
	clientHeightTotal.value += first.offsetHeight

	if (preListAfter.value.length <= 0) {
		return
	}
	// 更新数据
	preListBefore.value.push(dataList.value.shift())
	dataList.value.push(preListAfter.value.shift())
}

async function updateListBefore() {
	if (preListBefore.value.length <= 0) {
		return
	}

	// 更新数据
	dataList.value.unshift(preListBefore.value.pop())
	preListAfter.value.push(dataList.value.pop())

	await nextTick()

	const first = document.querySelector("#list0")

	// 计算元素移动的距离
	clientHeightTotal.value -= first.offsetHeight
	if (clientHeightTotal.value <= 0) {
		clientHeightTotal.value = 0
	}

	// curPage.value = Math.ceil(preListBefore.value.length / 10) + 1

}






// 加载页面获取数据
onMounted(async () => {
	const result = await getData(curPage.value)

	dataList.value = result.data
	total.value = result.total

	getPreData()
})



// 监听 preListAfter 数据是否用完， 用完就请求
watch(preListAfter, async (newValue) => {
	if (newValue.length == 0) {
		curPage.value++
		await getPreData()
		updateListAfter()
	}
}, { deep: true })



// 获取预览数据
async function getPreData() {
	if (curPage.value >= total.value) {
		return
	}

	const result = await getData(curPage.value + 1)

	preListAfter.value = result.data
}

// 获取数据
async function getData(cur = 1) {
	const result = await fetch(`http://127.0.0.1:8080/getPage?curPage=${cur}`).then(r => r.json())

	if (result.code != 200) {
		alert(result.msg)
		return
	}

	return result.data
}

</script>

<style scoped>
p {
	/* transform: translateY(100px); */
}

.list {
	height: 500px;
	overflow-y: auto;
	border: 5px solid aqua;
}

li {
	border: 1px dashed #ccc;
	margin: 10px 0;
}
</style>
