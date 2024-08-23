package api

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	"sublink/models"
	"sublink/node"
	"time"

	"github.com/gin-gonic/gin"
)

// NodeUpdate 节点更新
func NodeUpdate(c *gin.Context) {
	var nodeUpdating models.Node
	idStr := c.PostForm("id")
	name := c.PostForm("name")
	link := c.PostForm("link")
	if name == "" || link == "" {
		c.JSON(400, gin.H{
			"msg": "节点名称 or 备注不能为空",
		})
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// 处理转换错误，比如返回一个错误响应
		c.JSON(400, gin.H{
			"error": "无效的ID",
		})
		return
	}
	nodeUpdating.ID = id
	nodeUpdating.Name = name
	nodeUpdating.Link = link
	err = nodeUpdating.UpdateById()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "更新失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "更新成功",
	})
}

// NodeGet 获取节点列表
func NodeGet(c *gin.Context) {
	var Node models.Node
	nodes, err := Node.List()
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "node list error",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": nodes,
		"msg":  "node get",
	})
}

// NodeAdd 添加节点
func NodeAdd(c *gin.Context) {
	var Node models.Node
	link := c.PostForm("link")
	name := c.PostForm("name")
	if link == "" {
		c.JSON(400, gin.H{
			"msg": "link  不能为空",
		})
		return
	}
	if !strings.Contains(link, "://") {
		c.JSON(400, gin.H{
			"msg": "link 必须包含 ://",
		})
		return
	}
	Node.Name = name
	if name == "" {
		u, err := url.Parse(link)
		if err != nil {
			log.Println(err)
			return
		}
		switch {
		case u.Scheme == "ss":
			ss, err := node.DecodeSSURL(link)
			if err != nil {
				log.Println(err)
				return
			}
			Node.Name = ss.Name
		case u.Scheme == "ssr":
			ssr, err := node.DecodeSSRURL(link)
			if err != nil {
				log.Println(err)
				return
			}
			Node.Name = ssr.Qurey.Remarks
		case u.Scheme == "trojan":
			trojan, err := node.DecodeTrojanURL(link)
			if err != nil {
				log.Println(err)
				return
			}
			Node.Name = trojan.Name
		case u.Scheme == "vmess":
			vmess, err := node.DecodeVMESSURL(link)
			if err != nil {
				log.Println(err)
				return
			}
			Node.Name = vmess.Ps
		case u.Scheme == "vless":
			vless, err := node.DecodeVLESSURL(link)
			if err != nil {
				log.Println(err)
				return
			}
			Node.Name = vless.Name
		case u.Scheme == "hy" || u.Scheme == "hysteria":
			hy, err := node.DecodeHYURL(link)
			if err != nil {
				log.Println(err)
				return
			}
			Node.Name = hy.Name
		case u.Scheme == "hy2" || u.Scheme == "hysteria2":
			hy2, err := node.DecodeHY2URL(link)
			if err != nil {
				log.Println(err)
				return
			}
			Node.Name = hy2.Name
		case u.Scheme == "tuic":
			tuic, err := node.DecodeTuicURL(link)
			if err != nil {
				log.Println(err)
				return
			}
			Node.Name = tuic.Name
		}
	}
	Node.Link = link
	Node.CreateDate = time.Now().Format("2006-01-02 15:04:05")
	err := Node.Find()
	// 如果找到记录说明重复
	if err == nil {
		Node.Name = name + " " + time.Now().Format("2006-01-02 15:04:05")
	}
	err = Node.Add()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "添加失败检查一下是否节点重复",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "添加成功",
	})
}

// NodeDel 删除节点
func NodeDel(c *gin.Context) {
	var Node models.Node
	id := c.Query("id")
	if id == "" {
		c.JSON(400, gin.H{
			"msg": "id 不能为空",
		})
		return
	}
	x, _ := strconv.Atoi(id)
	Node.ID = x
	err := Node.Del()
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "删除失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"msg":  "删除成功",
	})
}

// NodesTotal 节点统计
func NodesTotal(c *gin.Context) {
	var Node models.Node
	nodes, err := Node.List()
	count := len(nodes)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "获取不到节点统计",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "00000",
		"data": count,
		"msg":  "取得节点统计",
	})
}

func NodeSort(c *gin.Context) {
	var nodes []models.Node
	// 解析请求体中的 JSON 数据
	if err := c.BindJSON(&nodes); err != nil {
		c.JSON(400, gin.H{"msg": "解析失败", "error": err.Error()})
		return
	}

	// 更新节点排序
	if err := models.UpdateSort(nodes); err != nil {
		c.JSON(500, gin.H{"msg": "排序更新失败", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code": "00000",
		"data": nil,
		"msg":  "排序更新成功",
	})
}
