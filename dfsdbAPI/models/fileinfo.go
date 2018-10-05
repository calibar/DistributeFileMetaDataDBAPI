package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Fileinfo struct {
	Id           int       `orm:"column(id);auto"`
	Guid         string    `orm:"column(guid)"`
	OriginalName string    `orm:"column(originalName);size(100)"`
	FileLocation string    `orm:"column(fileLocation)"`
	Timestamp    time.Time `orm:"column(timestamp);type(timestamp);auto_now"`
	ApplicationID string    `orm:"column(applicationID);size(255)"`
	ApplicationMetaData string    `orm:"column(applicationMetaData)"`
}

func (t *Fileinfo) TableName() string {
	return "fileinfo"
}

func init() {
	orm.RegisterModel(new(Fileinfo))
}

// AddFileinfo insert a new Fileinfo into database and returns
// last inserted Id on success.
func AddFileinfo(m *Fileinfo) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetFileinfoById retrieves Fileinfo by Id. Returns error if
// Id doesn't exist
func GetFileinfoById(id int) (v *Fileinfo, err error) {
	o := orm.NewOrm()
	v = &Fileinfo{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllFileinfo retrieves all Fileinfo matches certain condition. Returns empty list if
// no records exist
func sliceDelete(origin []Fileinfo,fileinfo Fileinfo)([]Fileinfo,int){
	var index int
	for i := 0; i < len(origin); i++ {
		if origin[i]==fileinfo {
			origin = append(origin[:i], origin[i+1:]...)
			i-- // maintain the correct index
			index=i
		}
	}
	return origin,index
}
func changeTimeFormat(t time.Time)time.Time{
	tstr:=t.Format("2006-01-02 15:04:05")
	t,_=time.Parse("2006-01-02 15:04:05",tstr)
	return t
}
func GetAllFileinfo(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64,timeFilter map[string]string) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Fileinfo))

	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
			fmt.Println(k)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}


	var l []Fileinfo
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		var beginTime time.Time
		var stopTime time.Time
		for k, v := range timeFilter {
			fmt.Println("k:",k,"v:",v)
			if k==""{
				var nl []Fileinfo
				stopTime,_=time.Parse("2006-01-02 15:04:05",v)
				fmt.Println(stopTime)
				for _, vl := range l {
					vlt:=changeTimeFormat(vl.Timestamp)
					if vlt.Before(stopTime){
						nl=append(nl,vl)
					}
				}
				l=nl
			}else if v==""{
				var nl []Fileinfo
				beginTime,_=time.Parse("2006-01-02 15:04:05",k)
				fmt.Println(beginTime)
				for _, vl := range l {
					vlt:=changeTimeFormat(vl.Timestamp)
					if vlt.After(beginTime){
						nl=append(nl,vl)
					}
				}
				l=nl
			}else {
				var nl []Fileinfo
				beginTime,_=time.Parse("2006-01-02 15:04:05",k)
				stopTime,_=time.Parse("2006-01-02 15:04:05",v)
				for _, vl := range l {
					vlt:=changeTimeFormat(vl.Timestamp)
					if vlt.After(beginTime)&&vlt.Before(stopTime){
						nl=append(nl,vl)
					}
				}
				l=nl
			}
			// rewrite dot-notation to Object__Attribute
			/*beginTime,_:=time.Parse("2018-10-05 19:00:00",k)
			stopTime,_:=time.Parse("2018-10-05 19:00:00",v)
			for _, v := range l {
				if v.Timestamp<=
			}*/
		}
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateFileinfo updates Fileinfo by Id and returns error if
// the record to be updated doesn't exist
func UpdateFileinfoById(m *Fileinfo) (err error) {
	o := orm.NewOrm()
	v := Fileinfo{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteFileinfo deletes Fileinfo by Id and returns error if
// the record to be deleted doesn't exist
func DeleteFileinfo(id int) (err error) {
	o := orm.NewOrm()
	v := Fileinfo{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Fileinfo{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
