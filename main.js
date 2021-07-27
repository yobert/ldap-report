function main() {
	let root = {"children":{}}
	let id = 0

	for(let [ groupkey, group ] of Object.entries(main_data.groups)) {
		let tree = root
		let treedn = ""

		for(let i = 0; i < group.parsed.length; i++) {
			let attr = group.parsed[i]
			let subdn = attr.type + "=" + attr.value // TODO implement ldap.EscapeFilter
			let key = (attr.value + "=" + attr.type).toLowerCase()

			if(treedn == "") {
				treedn = subdn
			} else {
				treedn = subdn + "," + treedn
			}

			if(!tree.children[key]) {
				let node = {
					id: "node"+id,
					dn: treedn,
					subdn: subdn,
					name: attr.value,
					children: {},
				}
				id++
				tree.children[key] = node
			}
			tree = tree.children[key]
		}

		tree.name = String.fromCodePoint(0x01f465) + " " + tree.name

		if(group.users) {
			for(let i = 0; i < group.users.length; i++) {
				let userkey = group.users[i]
				let user = main_data.users[userkey]
				let key = user.name.toLowerCase()

				tree.children[key] = {
					id: "node"+id,
					dn: user.dn,
					subdn: user.subdn,
					name: String.fromCodePoint(0x01f464) + " " + user.name,
					children: {},
				}
				id++
			}
		}
	}

	for(let [ userkey, user ] of Object.entries(main_data.users)) {
		let tree = root
		let treedn = ""

		for(let i = 0; i < user.parsed.length; i++) {
			let attr = user.parsed[i]
			let subdn = attr.type + "=" + attr.value // TODO implement ldap.EscapeFilter
			let key = (attr.value + "=" + attr.type).toLowerCase()

			if(treedn == "") {
				treedn = subdn
			} else {
				treedn = subdn + "," + treedn
			}

			if(!tree.children[key]) {
				let node = {
					id: "node"+id,
					dn: treedn,
					subdn: subdn,
					name: attr.value,
					children: {},
				}
				id++
				tree.children[key] = node
			}
			tree = tree.children[key]
		}

		tree.name = String.fromCodePoint(0x01f464) + " " + tree.name

		if(user.groups) {
			for(let i = 0; i < user.groups.length; i++) {
				let groupkey = user.groups[i]
				let group = main_data.groups[groupkey]
				let key = group.name.toLowerCase()

				tree.children[key] = {
					id: "node"+id,
					dn: group.dn,
					subdn: group.subdn,
					name: String.fromCodePoint(0x01f465) + " " + group.name,
					children: {},
				}
				id++
			}
		}
	}

	sortdata(root)

	let body = document.documentElement

	body.appendChild(build(root))

	//let pre = document.createElement("pre")
	//pre.appendChild(document.createTextNode(JSON.stringify(root, null, "    ")))
	//body.appendChild(pre)
}

function sortdata(node) {
	let keys = []
	for(let [ key, child ] of Object.entries(node.children)) {
		keys.push(key)
		sortdata(child)
	}
	keys.sort()
	let children = []
	for(let i = 0; i < keys.length; i++) {
		children.push(node.children[keys[i]])
	}
	node.children = children
	if(children.length > 5) {
		node.open = false
	} else {
		node.open = true
	}
}

function toggle(node) {
	let plus = document.getElementById("plus_"+node.id)
	let children = document.getElementById("children_"+node.id)

	if(node.open) {
		node.open = false
		plus.innerHTML = "+ "
		children.className = "children children_closed"
	} else {
		node.open = true
		plus.innerHTML = "- "
		children.className = "children"
	}
}

function build(node) {
	if(!node.dn) {
		// root node
		let div = document.createElement("div")
		for(let i = 0; i < node.children.length; i++) {
			div.appendChild(build(node.children[i]))
		}
		return div
	}

	let len = node.children.length

	let div = document.createElement("div")
	div.className = "node"

	let toggler;
	if(len > 0) {
		toggler = document.createElement("a")
		toggler.setAttribute("href", "#")
		toggler.addEventListener("click", function(e) {
			e.preventDefault()
			toggle(node)
			return false
		}, false)
	} else {
		toggler = document.createElement("span")
	}

	toggler.className = "name_toggle"

	toggler.setAttribute("title", node.dn)

	let plus = document.createElement("span")
	plus.className = "plus"
	if(len > 0) {
		if(node.open) {
			plus.appendChild(document.createTextNode("- "))
		} else {
			plus.appendChild(document.createTextNode("+ "))
		}
	} else {
		plus.appendChild(document.createTextNode("\u00a0\u00a0"))
	}

	plus.setAttribute("id", "plus_"+node.id)

	toggler.appendChild(plus)

	let name = document.createElement("span")
	name.className = "name"
	name.appendChild(document.createTextNode(node.name))

	toggler.appendChild(name)

	if(len > 0) {
		let count = document.createElement("span")
		count.className = "count"
		count.appendChild(document.createTextNode(len))
		toggler.appendChild(document.createTextNode(" "))
		toggler.appendChild(count)
	}

	div.appendChild(toggler)

	if(len == 0)
		return div

	let children = document.createElement("div")
	if(node.open)
		children.className = "children"
	else
		children.className = "children children_closed"

	children.setAttribute("id", "children_"+node.id)

	for(let i = 0; i < len; i++) {
		children.appendChild(build(node.children[i]))
	}

	div.appendChild(children)
	return div
}
