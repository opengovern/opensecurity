


import { Badge, FileUpload, FormField, Input } from "@cloudscape-design/components";
import {  CredentialField, Schema } from "./types";
import { dateTimeDisplay } from "../../../utilities/dateDisplay";

// For whole schema
// type 1 credential
// type 0 integration


export const GetActions=(type: number,schema: Schema | undefined)=>{
    if(type===1){
        return schema?.actions?.credentials
    }

    return schema?.actions?.integrations

}

export const GetTableColumns=(type: number,schema: Schema | undefined)=>{
    if(type===1){
        const fields = schema?.render?.credentials?.fields
        
        return fields
            ?.filter((field) => field?.show)
            ?.sort((a, b) => a.order - b.order)
            .map((field) => {
                return {
                    title: field.label,
                    dataIndex: field.name,
                    key: field.name,
                    sorter: field.sortable,
                    filter: field.filterable,
                    statusOptions: field.statusOptions,
                    type: field.fieldType,
                }
            })
        
    }

    const fields = schema?.render?.integrations?.fields
  

    return fields
        // ?.filter((field) => field?.show)
        ?.sort((a, b) => a.order - b.order)
        .map((field) => {
            return {
                title: field.label,
                dataIndex: field.name,
                key: field.name,
                sorter: field.sortable,
                filter: field.filterable,
                type: field.fieldType,
                statusOptions: field.statusOptions,
            }
        })
}

export const GetTableColumnsDefintion=(type: number,schema: Schema | undefined)=>{
   if (type === 1) {
       const fields = schema?.render.credentials.fields
       return fields
           ?.sort((a, b) => a.order - b.order)
           .map((field) => {
               return {
                   id: field.name,
                  visible: true,
               }
           })
   }

   const fields = schema?.render?.integrations?.fields
   return fields
       ?.sort((a, b) => a.order - b.order)
       .map((field) => {
           return {
               id: field.name,
               visible: true,
           }
       })
}

export const GetView=(type: number,schema: Schema)=>{
    if(type===1){
        return schema.render.credentials
    }

    return schema.render.integrations
}

export const GetDetails=(type: number,schema: Schema)=>{
    if(type===1){
        return schema.render.credentials
    }

    return schema.render.integrations
}

export const GetDetailsFields=(type: number,schema: Schema)=>{
    if(type===1){
        return schema.render.credentials.fields?.filter((field) => field.detail).sort(
            (a, b) => a.order - b.order
        )
    }

    return schema.render.integrations.fields
        .filter((field) => field.detail)
        .sort((a, b) => a.order - b.order)
}

export const GetDetailsActions=(type: number,schema: Schema |undefined)=>{
    if(type===1){
        return schema?.actions.credentials.map((action)=>{
            return action.type})
    }

    return schema?.actions.integrations.map((action) => {
        return action.type
    })
}

export const GetDefaultPageSize=(type: number,schema: Schema)=>{
    if(type===1){
        return schema.render.credentials.defaultPageSize
    }

    return schema.render.integrations.defaultPageSize
}

export const GetLogo=(schema: Schema)=>{
    return schema.icon
}

export const GetDiscover=(schema : Schema | undefined)=>{
    return schema?.discover.credentials.sort((a,b)=>a.priority-b.priority)
}

export const GetDiscoverField = (schema: Schema | undefined, index: number) => {
    return schema?.discover.credentials[index].fields.sort(
        (a, b) => a.order - b.order
    )
}
export const GetUpdateCredentialFields = (schema: Schema | undefined, index:number) => {
    const fields = schema?.discover.credentials[index].fields.sort(
        (a, b) => a.order - b.order
    )
    const editableFields = schema?.actions.credentials.filter((item: any)=>{
        return item.type === 'update'
    })[0].editableFields
    return fields?.filter((field) => editableFields?.includes(field.name))

}
export const GetEditField = (schema: Schema | undefined, type: number) => {
    if(type ==1){
        const actions = schema?.actions.credentials
        const edit_action = actions?.find((action) => action.type === 'update')
        const fields = edit_action?.editableFields
        return schema?.render?.credentials?.fields.filter((field) =>
            fields?.includes(field.name)
        )

    }
    const actions = schema?.actions.integrations
    const edit_action = actions?.find((action) => action.type === 'update')
    const fields = edit_action?.editableFields
    return schema?.render?.integrations?.fields.filter((field) =>
        fields?.includes(field.name)
    )


}

export const CheckCondition=(regex: string|undefined, value: any)=>{
    if(!value || value === ''){
        return true
    }
    if(regex){
    return new RegExp(regex).test(value)

    }
    else{
        return false
    }

}
export const CheckFileSize=(file: File, size: number)=>{
    if(file){
      


    return file?.size/1024 <= size

    }
    else{
        return true
    }
}
 function getBase64(file:any) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = (error) => reject(error);
  });
}
export const RenderInputField = (
    field: CredentialField,
    setFunction: Function,
    value: any
) => {
    // handle tet and password
    if (field.inputType == 'text' || field.inputType == 'password') {
        return (
            <FormField
                className="w-full"
                label={field.label}
                errorText={
                    !CheckCondition(field?.validation?.pattern, value)
                        ? field?.validation?.errorMessage
                        : ''
                }
            >
                <Input
                    className="w-full"
                    value={value}
                    type={field.inputType}
                    onChange={({ detail }) => setFunction(detail.value)}
                />
            </FormField>
        )
    }
    // handle file
    if (field.inputType == 'file') {
        return (
            <FormField className="w-full" label={field.label}>
                <FileUpload
                    className="w-full"
                    onChange={({ detail }) => setFunction(detail.value[0])}
                    value={value ? [value] : []}
                    i18nStrings={{
                        uploadButtonText: (e) =>
                            e ? 'Choose files' : 'Choose file',
                        dropzoneText: (e) =>
                            e ? 'Drop files to upload' : 'Drop file to upload',
                        removeFileAriaLabel: (e) => `Remove file ${e + 1}`,
                        limitShowFewer: 'Show fewer files',
                        limitShowMore: 'Show more files',
                        errorIconAriaLabel: 'Error',
                    }}
                    showFileLastModified
                    showFileSize
                    multiple
                    showFileThumbnail
                    tokenLimit={3}
                    accept={`${field.validation.fileTypes?.map((item) => {
                        return `${item},`
                    })}`}
                    errorText={
                        CheckFileSize(
                            value,
                            field.validation.maxFileSizeMB ?? 0
                        )
                            ? ''
                            : field.validation.errorMessage
                    }
                />
            </FormField>
        )
    }
}


export const RenderUpdateField = (
    field: CredentialField,
    setFunction: Function,
    value: any
) => {
    // handle tet and password
    if (field.inputType == 'text' || field.inputType == 'password') {
        return (
            <FormField
                className="w-full"
                label={field.label}
                errorText={
                    !CheckCondition(field.validation.pattern, value)
                        ? field.validation.errorMessage
                        : ''
                }
            >
                <Input
                    className="w-full"
                    value={value}
                    type={field.inputType}
                    onChange={({ detail }) => setFunction(detail.value)}
                />
            </FormField>
        )
    }
    // handle file
    if (field.inputType == 'file') {
        return (
            <FormField
                className="w-full"
                label={field.label}
                errorText={
                    !CheckCondition(field.validation.pattern, value)
                        ? field.validation.errorMessage
                        : ''
                }
            >
                <FileUpload
                    className="w-full"
                    onChange={({ detail }) => setFunction(detail.value)}
                    value={value ? [value] : []}
                    i18nStrings={{
                        uploadButtonText: (e) =>
                            e ? 'Choose files' : 'Choose file',
                        dropzoneText: (e) =>
                            e ? 'Drop files to upload' : 'Drop file to upload',
                        removeFileAriaLabel: (e) => `Remove file ${e + 1}`,
                        limitShowFewer: 'Show fewer files',
                        limitShowMore: 'Show more files',
                        errorIconAriaLabel: 'Error',
                    }}
                    showFileLastModified
                    showFileSize
                    accept={`${field.validation.fileTypes
                        ?.map((type) => `${type}`)
                        .join(',')}`}
                    showFileThumbnail
                    tokenLimit={1}
                    errorText={
                        CheckFileSize(
                            value,
                            field.validation.maxFileSizeMB ?? 0
                        )
                            ? ''
                            : field.validation.errorMessage
                    }
                />
            </FormField>
        )
    }
}

export const RenderTableField = (field: any, item: any) => {
    if (field.type === 'text') {
        return item[field.key]
    }
    if (field.type === 'status') {
     
    
        return (
            <Badge
                // @ts-ignore
                color={
                    // @ts-ignore
                    field?.statusOptions?.filter((x: any) => {
                        
                        if (x.value === item[field.key]) {
                            return x
                        }
                    })[0]?.color ?? 'green'
                }
                // @ts-ignore
            >
                {/* @ts-ignore */}
                {item[field.key]}
            </Badge>
        )
    }
    if (field.type === 'date') {
        return dateTimeDisplay(item[field.key])
    }
    return item[field.key]

}

export const GetViewFields = (schema: Schema | undefined, type: number) => {
    if (type === 1) {
        return schema?.render?.credentials?.fields
            .filter((field) => {
                if(field.detail){
                    return field
                }
            })
            .sort((a, b) => a.order - b.order)
            ?.map((field) => {
                return {
                    title: field.label,
                    dataIndex: field.name,
                    key: field.name,
                    type: field.fieldType,
                    statusOptions: field.statusOptions,
                }
            })
    }
    return schema?.render?.integrations?.fields
        .filter((field) => {
            if (field.detail) {
                return field
            }
        })
        .sort((a, b) => a.order - b.order)
        ?.map((field) => {
            return {
                title: field.label,
                dataIndex: field.name,
                key: field.name,
                type: field.fieldType,
                statusOptions: field.statusOptions,
            }
        })
}

