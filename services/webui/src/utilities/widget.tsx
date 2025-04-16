
export interface Widget {
    id:           string;
    data:         Data;
    rowSpan:      number;
    columnSpan:   number;
    columnOffset: number;
}

export interface ColumnOffset {
    "4": number;
}

export interface Data {
    componentId: string
    title: string
    description: string
    user_id: string
    is_public: boolean
    props: any
}



export interface Kpi {
    info:      string;
    count_kpi: string;
    list_kpi:  string;
}

export interface WidgetAPI {
    id:           string;
    title:        string;
    description: string;
    widget_type: string;
    widget_props: any;
    user_id:   string;
    is_public: boolean;
    row_span:      number;
    column_span:   number;
    column_offset: number;
}

export const  WidgetToAPI = (widget: Widget,user_id : string,is_public: boolean) => {
    const widgetAPI  = {} as WidgetAPI
    widgetAPI.id = widget.id
    widgetAPI.title = widget.data.title
    widgetAPI.description = widget.data.description
    widgetAPI.widget_type = widget.data.componentId
    widgetAPI.widget_props = widget.data.props
    widgetAPI.user_id = user_id
    widgetAPI.is_public = is_public
    widgetAPI.row_span = widget.rowSpan
    widgetAPI.column_span = widget.columnSpan
    widgetAPI.column_offset = widget.columnOffset
    return widgetAPI

}
export const APIToWidget = (widget: WidgetAPI) => {
    const widgetData = {} as Widget
    widgetData.id = widget.id
    widgetData.data = {
        componentId: widget.widget_type,
        title: widget.title,
        description: widget.description,
        user_id: widget.user_id,
        is_public: widget.is_public,
        props: widget.widget_props,
    }
    widgetData.rowSpan = widget.row_span
    widgetData.columnSpan = widget.column_span
    widgetData.columnOffset = widget.column_offset
    return widgetData
}